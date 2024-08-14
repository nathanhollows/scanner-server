package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/charmbracelet/log"
)

type trip struct {
	ID          string
	Name        sql.NullString
	Description sql.NullString
	Link        sql.NullString
}

func tripHandler(w http.ResponseWriter, r *http.Request) {
	switch r.PathValue("id") {
	case "":
		w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>My trip</title>
		</head>
		<body>

		<h1>My trip</h1>

		<form method="post">
			<label for="trip">Trip:</label><br>
			<input type="text" id="trip" name="trip"><br>	
			<label for="identifier">Identifier:</label><br>
			<input type="text" id="identifier" name="identifier"><br>
			<input type="submit" value="Submit">
		</form>

		</body>
		</html>
		`))
	default:
		tag := r.PathValue("id")

		rows, err := db.Query(`
			SELECT locations.*
			FROM locations 
			INNER JOIN scans ON locations.id = scans.location_id
			WHERE scans.tag = ?
			GROUP BY locations.id
			ORDER BY scans.timestamp DESC`, tag)
		if err != nil {
			log.Error("error querying database", "err", err)
			http.Error(w, "error querying database", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		trips := []trip{}
		for rows.Next() {
			scan := trip{}
			err = rows.Scan(&scan.ID, &scan.Name, &scan.Description, &scan.Link)
			if err != nil {
				log.Error("error scanning row", "err", err)
				http.Error(w, "error scanning row", http.StatusInternalServerError)
				return
			}
			trips = append(trips, scan)
		}
		json.NewEncoder(w).Encode(trips)

	}
}
