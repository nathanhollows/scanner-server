package handlers

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*
var templateFS embed.FS

func Render(
	w http.ResponseWriter,
	r *http.Request,
	tmpl string,
	block string,
	data interface{},
) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFS(templateFS, "templates/"+tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, block, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
