package handlers

import (
	"database/sql"
	"net/http"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

const dbname = "scans.db"

var db *sql.DB

func Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexAction)
	mux.HandleFunc("/register", registerAction)
	mux.HandleFunc("/scan", scanAction)
	mux.HandleFunc("/trip/", tripHandler)
	mux.HandleFunc("/trip/{id}", tripHandler)
	mux.HandleFunc("/generate", generateHandler)

	var err error
	db, err = sql.Open("sqlite3", dbname)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS scans (location_id TEXT, tag TEXT KEY, timestamp TEXT KEY)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS nodes (node TEXT PRIMARY KEY, location_id TEXT KEY)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS locations 
		(id INTEGER PRIMARY KEY, name TEXT KEY, description TEXT KEY, link TEXT KEY)`)
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":8088", mux))
}
