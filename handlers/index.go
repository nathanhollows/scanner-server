package handlers

import "net/http"

func indexAction(w http.ResponseWriter, r *http.Request) {
	Render(w, r, "index.html", "index", nil)
}
