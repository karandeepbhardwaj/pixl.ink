package server

import (
	"net/http"
	"time"

	"github.com/karandeepbhardwaj/pixl.ink/internal/config"
	"github.com/karandeepbhardwaj/pixl.ink/internal/handler"
)

func New(cfg *config.Config, h *handler.Handler) http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /{$}", h.Index)
	mux.HandleFunc("POST /upload", h.Upload)
	mux.HandleFunc("GET /qr/{id}", h.QRCode)
	mux.HandleFunc("GET /{id}", h.View)

	mux.HandleFunc("POST /api/upload", h.APIUpload)
	mux.HandleFunc("GET /api/health", h.Health)

	rl := NewRateLimiter(30, time.Minute)
	var chain http.Handler = mux
	chain = rl.Middleware(chain)
	chain = LoggingMiddleware(chain)
	chain = RecoveryMiddleware(chain)

	return chain
}
