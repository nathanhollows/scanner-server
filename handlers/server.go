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
	mux.HandleFunc("/scan", scanAction)
	mux.HandleFunc("/generate", generateHandler)

	log.Fatal(http.ListenAndServe(":8088", mux))
}
