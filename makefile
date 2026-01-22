TARGETDIR=./deploy
proj=github.com/alexandr-andreyev/soup-rk7-events

# Version info
VERSION ?= 1.0.0
COMMIT := $(shell git rev-parse --short HEAD 2>NUL || echo unknown)
BUILD_TIME := $(shell powershell -Command "Get-Date -Format 'yyyy-MM-dd HH:mm:ss'" 2>NUL || date +"%Y-%m-%d %H:%M:%S")

LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X 'main.buildTime=$(BUILD_TIME)'

all: vet test buildEXE

vet:
	go vet ./cmd/app
	go vet ./internal/...

test:
	go test -timeout 30s ./internal/...

buildEXE:
	go build -o "./souprk7notify.exe" -a -ldflags "$(LDFLAGS)" ./cmd/app

run:
	go run ./cmd/app debug

version:
	@echo Version: $(VERSION)
	@echo Commit: $(COMMIT)
	@echo Build time: $(BUILD_TIME)