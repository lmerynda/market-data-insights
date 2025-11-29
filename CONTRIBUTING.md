# Development

## Prerequisites

- Go 1.23.8 or later

## Run Locally

### Option 1: Direct Go
```bash
go mod download
go run main.go
```

### Option 2: Build Script
```bash
./run.sh
```

### Option 3: VS Code Dev Container
Open in VS Code and select "Reopen in Container" when prompted.

Application runs on http://localhost:8080

## Build

```bash
go build -o market-data-insights main.go
./market-data-insights
```