TARGETDIR=./deploy
proj=github.com/alexandr-andreyev/soup-rk7-events

# Get git commit hash (use 'unknown' if git command fails on Windows)
sha1ver := $(shell git rev-parse HEAD || echo unknown)

all: vet test buildEXE

vet:
	go vet ./cmd/app
	go vet ./internal/...

test:
	go test -timeout 30s ./internal/...

buildEXE:
	go build -o "./souprk7notify.exe" -a -ldflags "-X main.sha1ver=$(sha1ver)" ./cmd/app

run:
	go run ./cmd/app debug