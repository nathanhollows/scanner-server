package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
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

		if node == "" || tag == "" {
			log.Info("malformed request to /scan: node and tag must be provided")
			http.Error(w, "node and tag must be provided", http.StatusBadRequest)
			return
		}

		// Insert the scan into the database
		result, err := db.Exec(`
		INSERT INTO scans (location_id, tag, timestamp) 
		SELECT locations.id, ?, datetime('now') 
			FROM nodes
		INNER JOIN locations ON nodes.location_id = locations.id
		WHERE nodes.node = ?`, tag, node)
		if err != nil {
			log.Error("error inserting scan", "err", err)
			http.Error(w, "error inserting scan into database", http.StatusInternalServerError)
			return
		}
		// Check how many rows were affected
		rows, _ := result.RowsAffected()
		if rows == 0 {
			log.Error("node was not found or has not been registered", "node", node)
			http.Error(w, "node was not found or has not been registered", http.StatusNotFound)
			return
		}

		log.Info("scan recorded", "node", node, "tag", tag)
		http.Error(w, "scan recorded", http.StatusOK)
		return

	default:
		log.Info("method not allowed on /scan", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
