# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Windows service that receives XML-formatted events from R_Keeper 7 (RK7) POS system via HTTP notify interface. Events are processed and forwarded to external systems (e.g., stop-list changes, order modifications).

## Common Commands

### Build
```bash
make buildEXE          # Build souprk7notify.exe
make all              # Run vet, test, and build
make run              # Run in debug mode
go run .\cmd\app debug  # Alternative debug mode
```

### Testing & Quality
```bash
make vet              # Run go vet with -all -shadow
make test             # Run tests with 30s timeout
```

### Service Management (Admin Command Prompt)
```cmd
souprk7notify.exe install    # Install as Windows service
souprk7notify.exe remove     # Remove service
souprk7notify.exe start|stop|pause|continue
```

**Service Identity:**
- Name: `soup-rk7-events`
- Display: `Soup Events for R_Keeper7`

## Architecture

### Request Flow
```
RK7 Server → POST :7080/events (XML) → Handler → NotifyEventService → External Systems
```

### Component Structure

**Entry & Service Control** (`cmd/app/`):
- `main.go`: Service names, entry point
- `svc_service.go`: Windows service Execute() loop, handles stop/pause/continue
- `svc_*.go`: Installation, management, argument parsing
- Uses `golang.org/x/sys/windows/svc` for service integration

**Application Core** (`internal/app/`):
- `run.go`: Launches HTTP server as goroutine
- `setup.go`: Initializes NotifyEventService and HTTP server (:7080)
- `server.go`: Server struct holds winlog and httpServer

**HTTP Layer** (`internal/transport/`):
- `handler.go`: `/events` endpoint, unmarshals XML → calls NotifyEventService

**Business Logic** (`internal/services/`):
- `rkeeperNotifyHandleService.go`: Interface `NotifyEventService` with method `HandleNotification()`
- Current implementation: logs event name, placeholder for forwarding logic

**Data Models** (`internal/models/`):
- `rk7NotifyEvent.go`: XML structures for RK7 events
  - Root: `Rk7NotifyEvent` (RestCode, DateTime, Situation, Name, etc.)
  - Nested: Station, Server, Waiter, Order, Item (all optional pointers)

### Key Configuration Points

- **HTTP Port**: Hardcoded `:7080` in `internal/app/setup.go:46`
- **Timeouts**: 10s read/write in `internal/app/setup.go:48-49`
- **Service Names**: Hardcoded in `cmd/app/main.go:16-19` (TODO: move to config)

### Extending Event Handling

To handle new RK7 event types:
1. Add case in `internal/services/rkeeperNotifyHandleService.go:21-28` switch on `data.Name`
2. Access nested data from `data.Station`, `data.Order`, etc. as needed
3. Implement external forwarding logic

### Logging

- **Debug mode**: stdout via `debug.Log`
- **Service mode**: Windows Event Log via `eventlog`
- HTTP handler logs all incoming requests and parsed events

## Known Issues

- Makefile vet target references old `cmd/gosvc` path (should be `cmd/app`)
- SHA1 version embedding noted as "not working" in makefile
- Service names should come from configuration file