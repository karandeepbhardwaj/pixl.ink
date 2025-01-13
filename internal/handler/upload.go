package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/karandeepbhardwaj/pixl.ink/internal/config"
	"github.com/karandeepbhardwaj/pixl.ink/internal/image"
	"github.com/karandeepbhardwaj/pixl.ink/internal/shortid"
	"github.com/karandeepbhardwaj/pixl.ink/internal/storage"
)

type Handler struct {
	cfg  *config.Config
	db   *storage.SQLiteStore
	disk *storage.DiskStore
}

func New(cfg *config.Config, db *storage.SQLiteStore, disk *storage.DiskStore) *Handler {
	return &Handler{cfg: cfg, db: db, disk: disk}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxUploadSize)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := image.Validate(header, h.cfg.MaxUploadSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := shortid.Generate(6)
	if err != nil {
		http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
		return
	}

	ext := image.GetExtension(header.Header.Get("Content-Type"))
	filename := id + ext

	if _, err := h.disk.Save(filename, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	meta := &storage.ImageMeta{
		ShortID:     id,
		Filename:    filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
	}
	if err := h.db.Save(meta); err != nil {
		http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/%s", h.cfg.BaseURL, id)

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<div class="result-card">
			<img src="/uploads/%s" alt="Uploaded image" class="result-image">
			<div class="result-url">
				<input type="text" value="%s" readonly id="shareUrl">
				<button onclick="copyUrl()" class="btn-copy">Copy</button>
			</div>
			<img src="/qr/%s" alt="QR Code" class="result-qr">
		</div>`, filename, url, id)
		return
	}

	http.Redirect(w, r, "/"+id, http.StatusSeeOther)
}

func (h *Handler) APIUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxUploadSize)

	file, header, err := r.FormFile("file")
	if err != nil {
		jsonError(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := image.Validate(header, h.cfg.MaxUploadSize); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := shortid.Generate(6)
	if err != nil {
		jsonError(w, "Failed to generate ID", http.StatusInternalServerError)
		return
	}

	ext := image.GetExtension(header.Header.Get("Content-Type"))
	filename := id + ext

	if _, err := h.disk.Save(filename, file); err != nil {
		jsonError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	meta := &storage.ImageMeta{
		ShortID:     id,
		Filename:    filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
	}
	if err := h.db.Save(meta); err != nil {
		jsonError(w, "Failed to save metadata", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":     id,
		"url":    fmt.Sprintf("%s/%s", h.cfg.BaseURL, id),
		"qr_url": fmt.Sprintf("%s/qr/%s", h.cfg.BaseURL, id),
	})
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
