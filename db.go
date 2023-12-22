package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they do not exist
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS pixels (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid TEXT UNIQUE,
            title TEXT
        );
        CREATE TABLE IF NOT EXISTS stats (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            pixel_id INTEGER,
            view_time DATETIME DEFAULT CURRENT_TIMESTAMP,
            ip TEXT,
            user_agent TEXT,
			fingerprint TEXT,
            FOREIGN KEY(pixel_id) REFERENCES pixels(id)
        );
    `)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func addStat(pixelID int, ip, userAgent string, fingerprint string) error {
	statement, err := db.Prepare("INSERT INTO stats (pixel_id, ip, user_agent, fingerprint) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(pixelID, ip, userAgent, fingerprint)
	return err
}

func getPixelIDFromUUID(uuid string) (int, error) {
	var id int
	row := db.QueryRow("SELECT id FROM pixels WHERE uuid = ?", uuid)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func savePixel(title, uuid string) error {
	statement, err := db.Prepare("INSERT INTO pixels (title, uuid) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(title, uuid)
	return err
}
