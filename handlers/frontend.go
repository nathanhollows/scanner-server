package handlers

import (
	"html/template"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/scanner-server/markdown"
	"github.com/nathanhollows/scanner-server/models"
	"github.com/nathanhollows/scanner-server/templates"
)

func indexAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// Redirect if ID is present.
	r.ParseForm()
	if r.Form.Has("id") {
		contentHandler(w, r)
		return
	}

	// Render the index page.
	c := templates.Index()
	err := templates.Layout(c).Render(r.Context(), w)
	if err != nil {
		log.Error("Failed to render index", "error", err)
	}
}

func contentHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("User entered ID", "id", r.FormValue("id"))

	tag := r.FormValue("id")

	scans := models.FindScansByTag(r.Context(), tag)
	if len(scans) == 0 {
		log.Info("No scans found for tag", "tag", tag)
	}

	var content []template.HTML

	md, err := markdown.RenderFromFile("bonus")
	// Ignore error if bonus.md is not found.
	if err == nil {
		content = append(
			content,
			md,
		)
	}

	for _, scan := range scans {
		md, err := markdown.RenderFromFile(scan.LocationID)
		if err != nil {
			log.Error("Failed to render markdown", "error", err, "tag", tag, "file", scan.LocationID)
		} else {
			content = append(
				content,
				md,
			)
		}
	}

	c := templates.Content(tag, content)
	err = templates.Layout(c).Render(r.Context(), w)
	if err != nil {
		log.Error("Failed to render content", "error", err)
	}
}
