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
- `main.go`: Service names (`svcName`, `svcNameLong`), calls `app.Run()`
- `svc_service.go`: Windows service Execute() loop, handles stop/pause/continue signals
- `svc_main.go`: `runService()` sets up logging (debug vs eventlog)
- `svc_install.go`, `svc_manage.go`: Service installation and management
- Uses `golang.org/x/sys/windows/svc` for Windows service integration

**Application Core** (`internal/app/`):
- `run.go`: Entry point - calls `setup()` then launches `runApp()` as goroutine
- `setup.go`: Initializes `NotifyEventService` and HTTP server with mux routing
- `app.go`: `runApp()` starts HTTP server via `ListenAndServe()`
- `server.go`: Server struct holds `httpServer` reference

**HTTP Layer** (`internal/transport/`):
- `handler.go`: `/events` POST endpoint
  - Reads body, unmarshals XML → `models.Rk7NotifyEvent`
  - Calls `NotifyService.HandleNotification()`
  - Logs request body and parsed events

**Business Logic** (`internal/services/`):
- `rkeeperNotifyHandleService.go`:
  - Interface: `NotifyEventService` with `HandleNotification()` method
  - Implementation: `RkNotifyHandleService` - currently logs event name, placeholder for forwarding

**Data Models** (`internal/models/`):
- `rk7NotifyEvent.go`: XML structures for RK7 events
  - Root: `Rk7NotifyEvent` (RestCode, DateTime, Situation, Name, ShiftNum, etc.)
  - Nested: `Station`, `Server`, `Waiter` (with `Role`), `Order`, `Item` (all optional pointers)

### Key Configuration Points

- **HTTP Port**: Hardcoded `:7080` in `internal/app/setup.go:38`
- **Timeouts**: 10s read/write in `internal/app/setup.go:40-41`
- **Service Names**: Hardcoded in `cmd/app/main.go:16,19` (TODO: move to config)
- **HTTP Route**: `/events` handler registered in `internal/app/setup.go:35`

### Extending Event Handling

To handle new RK7 event types:
1. Add case in `internal/services/rkeeperNotifyHandleService.go:21` switch on `data.Name`
2. Access nested data from `data.Station`, `data.Order`, `data.Waiter`, etc.
3. Implement external forwarding logic (HTTP POST, message queue, etc.)

### Logging

- **Debug mode**: stdout via `debug.Log` (set in `svc_main.go:71`)
- **Service mode**: Windows Event Log via `eventlog` (set in `svc_main.go:73`)
- HTTP handler logs all incoming requests at `handler.go:24,34`

## Known Issues

- Makefile vet target line 10 references old `.\cmd\gosvc` path (should be `.\cmd\app`)
- Makefile test target line 14 references wrong module path (should use relative or correct module)
- SHA1 version embedding in makefile noted as "not working" (line 16)
- Service names should come from configuration file (TODOs in `main.go:16,19`)