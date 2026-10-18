# Market Data Insights

A Go web application for market data analytics, based on the market-watch-go project structure.

## Quick Start

Requires Go 1.26 or later.

```bash
go run main.go
```

Visit http://localhost:8080 to see the application.

## Features

- Simple web interface with Bootstrap styling
- RESTful API with health check endpoint
- Static file serving (CSS, JS)
- HTML templating with Gin

## API Endpoints

- `GET /` - Main dashboard page
- `GET /api/health` - API health check

## Project Structure

```
market-data-insights/
├── main.go              # Main application entry point
├── go.mod               # Go module dependencies
├── web/
│   ├── templates/       # HTML templates
│   │   └── index.html
│   └── static/          # Static assets
│       ├── css/
│       │   └── style.css
│       └── js/
│           └── app.js
└── README.md
```

## Dependencies

- Gin Web Framework
- Bootstrap 5 (CDN)

## Development

This is a minimal hello world web application with similar structure to market-watch-go but simplified for basic web functionality.
