package server

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var onePixelGIF = []byte{
	0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00,
	0x01, 0x00, 0x80, 0xff, 0x00, 0xff, 0xff, 0xff,
	0x00, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00, 0x00,
	0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44,
	0x01, 0x00, 0x3b,
}

func Server() {
	db = initDB("data.sqlite")

	http.HandleFunc("/new", newPixelPageHandler)
	http.HandleFunc("/create-pixel", createPixelHandler)
	http.HandleFunc("/p/", pixelHandler)
	http.HandleFunc("/dashboard", dashboardHandler)

	addr := ":8080"

	fmt.Printf("Server is up and running on http://localhost%s\n", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

type PixelStats struct {
	Title       string
	TotalViews  int
	UniqueViews int
	UUID        string
}
