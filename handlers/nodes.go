package handlers

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/scanner-server/models"
)

// linkAction links tags and cards
func linkAction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tagID := r.FormValue("tag")
		listID := r.FormValue("list")

		tag := models.Tag{
			TagID:  tagID,
			ListID: listID,
		}

		if tagID == "" || listID == "" {
			log.Info("malformed request to /link: tag and list must be provided")
			http.Error(w, "tag and list must be provided", http.StatusBadRequest)
			return
		}

		err := tag.Save(r.Context())
		if err != nil {
			log.Error("error inserting tag", "err", err)
			http.Error(w, "error inserting tag into database", http.StatusInternalServerError)
			return
		}

		log.Info("tag linked to list", "tag", tagID, "list", listID)
		w.Write([]byte("tag linked to list"))
	default:
		log.Info("method not allowed on /link", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

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
	case http.MethodGet:
		location := r.FormValue("location")
		tag := r.FormValue("tag")

		// Multiple scans are fine since we filter them later
		scan := models.Scan{
			LocationID: location,
			TagID:      tag,
			Timestamp:  time.Now(),
		}

		err := scan.Save(r.Context())
		if err != nil {
			log.Error("error inserting scan", "err", err)
			http.Error(w, "error inserting scan into database", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("scan recorded for location: " + location + " with tag: " + tag))
	default:
		log.Info("method not allowed on /scan", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
