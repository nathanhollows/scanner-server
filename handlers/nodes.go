package handlers

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/scanner-server/models"
)

// registerAction handles the node registration action.
// It accepts HTTP POST requests.
// For any other request method, it returns a 405 error.
// It implements the spec at
// https://github.com/nathanhollows/museum-scanner/blob/main/docs/api-spec.md
func registerAction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		node := r.FormValue("node")
		identifier := r.FormValue("identifier")

		if node == "" || identifier == "" {
			log.Info("malformed request to /register: node and identifier must be provided")
			http.Error(w, "node and identifier must be provided", http.StatusBadRequest)
			return
		}

		// Insert or update the node in the database
		_, err := db.Exec("INSERT INTO nodes (node, location_id) VALUES (?, ?)", node, identifier)
		if err.Error() == "UNIQUE constraint failed: nodes.node" {
			_, err := db.Exec("UPDATE nodes SET location_id = ? WHERE node = ?", identifier, node)
			if err != nil {
				log.Error("error updating node", "err", err)
				http.Error(w, "error updating node in database", http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			log.Error("error inserting node", "err", err)
			http.Error(w, "error inserting node into database", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO locations (id) VALUES (?) ON CONFLICT(id) DO NOTHING", identifier)
		if err != nil {
			log.Error("error inserting location", "err", err)
			http.Error(w, "error inserting location into database", http.StatusInternalServerError)
			return
		}

		log.Info("node registered", "node", node, "location", identifier)
		http.Error(w, "node registered", http.StatusOK)
		return

	default:
		log.Info("method not allowed on /register", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// scanAction handles users scanning a tag.
// It accepts HTTP POST requests.
// For any other request method, it returns a 405 error.
// It implements the spec at
// https://github.com/nathanhollows/museum-scanner/blob/main/docs/api-spec.md
func scanAction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		node := r.FormValue("node")
		tag := r.FormValue("tag")

		// Multiple scans are fine since we filter them later
		scan := models.Scan{
			LocationID: node,
			TagID:      tag,
			Timestamp:  time.Now(),
		}

		err := scan.Save(r.Context())
		if err != nil {
			log.Error("error inserting scan", "err", err)
			http.Error(w, "error inserting scan into database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	default:
		log.Info("method not allowed on /scan", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
