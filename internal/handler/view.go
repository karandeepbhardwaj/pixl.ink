package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	meta, err := h.db.GetByShortID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		filepath.Join("web", "templates", "layout.html"),
		filepath.Join("web", "templates", "view.html"),
	))

	data := map[string]interface{}{
		"Image":   meta,
		"URL":     h.cfg.BaseURL + "/" + meta.ShortID,
		"QRURL":   h.cfg.BaseURL + "/qr/" + meta.ShortID,
		"FileURL": "/uploads/" + meta.Filename,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "layout", data)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("web", "templates", "layout.html"),
		filepath.Join("web", "templates", "index.html"),
	))
	w.Header().Set("Content-Type", "text/html")
	tmpl.ExecuteTemplate(w, "layout", nil)
}
