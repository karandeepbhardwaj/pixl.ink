package handler

import (
	"encoding/json"
	"net/http"

	"github.com/karandeepbhardwaj/pixl.ink/internal/qr"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) QRCode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	meta, err := h.db.GetByShortID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	url := h.cfg.BaseURL + "/" + meta.ShortID
	png, err := qr.Generate(url, 256)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(png)
}
