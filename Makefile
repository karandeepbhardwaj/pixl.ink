.PHONY: build dev test clean

build:
	go build -o bin/pixlink main.go

dev:
	go run main.go

test:
	go test ./...

clean:
	rm -rf bin/ uploads/ *.db
