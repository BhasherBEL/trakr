package server

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func newPixelPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("public/html/new.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createPixelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	title := r.FormValue("title")

	pixelUUID := uuid.New().String()

	err := savePixel(title, pixelUUID)
	if err != nil {
		http.Error(w, "Error saving pixel to the database", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func pixelHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	uuid := pathParts[2]

	// Get the pixel ID from the database
	pixelID, err := getPixelIDFromUUID(uuid)
	if err != nil {
		log.Println("Error getting pixel ID:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	addr := strings.Split(r.RemoteAddr, ":")

	rawFingerprint := fmt.Sprintf(
		"%s|%s|%s",
		addr[0],
		r.Header.Get("User-Agent"),
		r.Header.Get("Accept-Language"),
	)
	hash := sha256.Sum256([]byte(rawFingerprint))
	fingerprint := string(hash[:])

	err = addStat(pixelID, addr[0], r.UserAgent(), fingerprint)
	if err != nil {
		log.Println("Error updating stats:", err)
	}

	// Serve the pixel image
	w.Header().Set("Content-Type", "image/gif")
	w.Write(onePixelGIF)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT p.title, p.uuid, COUNT(s.id) AS total_views, COUNT(DISTINCT s.ip) AS unique_views
		FROM pixels p
		LEFT JOIN stats s ON p.id = s.pixel_id
		GROUP BY p.id;`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		log.Println("Error fetching data:", err)
		return
	}
	defer rows.Close()

	var stats []PixelStats
	for rows.Next() {
		var s PixelStats
		if err := rows.Scan(&s.Title, &s.UUID, &s.TotalViews, &s.UniqueViews); err != nil {
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			log.Println("Error reading data:", err)
			return
		}
		stats = append(stats, s)
	}

	// Handle any errors encountered during iteration
	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating data", http.StatusInternalServerError)
		log.Println("Error iterating data:", err)
		return
	}

	// Render the template
	tmpl, err := template.ParseFiles("public/html/dashboard.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, stats)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
