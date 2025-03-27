# pixl.ink

[![CI](https://github.com/karandeepbhardwaj/pixl.ink/actions/workflows/ci.yml/badge.svg)](https://github.com/karandeepbhardwaj/pixl.ink/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Self-hosted image sharing with instant URLs and QR codes.

## Features

- Drag-and-drop image upload
- Instant shareable URLs with short IDs
- QR code generation for every image
- Dark/light mode
- JSON API for programmatic use
- Docker-ready deployment
- SQLite storage (zero external dependencies)
- Rate limiting and file size validation

## Quick Start

### Docker

```sh
docker compose up
```

### Local Development

```sh
go run main.go
```

Open [http://localhost:8080](http://localhost:8080)

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `MAX_UPLOAD_SIZE` | `10485760` | Max upload size in bytes (10MB) |
| `BASE_URL` | `http://localhost:8080` | Public base URL for generated links |
| `UPLOAD_DIR` | `./uploads` | Directory for uploaded files |
| `DB_PATH` | `./pixlink.db` | SQLite database path |

## API

### Upload Image

```sh
curl -X POST -F "file=@image.png" http://localhost:8080/api/upload
```

Response:

```json
{
  "id": "a7Bx3k",
  "url": "http://localhost:8080/a7Bx3k",
  "qr_url": "http://localhost:8080/qr/a7Bx3k"
}
```

### Health Check

```sh
curl http://localhost:8080/api/health
```

## Tech Stack

- **Go 1.22+** with stdlib HTTP server (enhanced ServeMux)
- **HTMX** for interactive uploads without page reloads
- **SQLite** via modernc.org/sqlite (pure Go, no CGO)
- **go-qrcode** for QR code generation

## License

MIT
