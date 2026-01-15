# URLShortener

A simple URL shortening REST API service built with Go. This project allows you to shorten long URLs, track clicks (planned, in pipeline), and easily redirect visitors using generated short links.

## Features

- Shorten long URLs to compact codes
- Redirect users from short codes to original URLs
- Click tracking for each shortened URL (planned, //todo)
- RESTful API endpoints
- Health check and database test endpoints
- Built with modular Go code and uses a PostgreSQL database

## Getting Started

### Prerequisites

- Go (recommended version 1.25+)
- Docker (optional, good to have, for containerized setup)
- PostgreSQL database

### Installation

#### 1. Clone the repository

```bash
git clone https://github.com/VxDant/URLShortener.git
cd URLShortener
```

#### 2. Set up the Database

You need a running PostgreSQL instance. Update the connection details as needed, or use the default provided in `scripts/migrate.sh`.

- Run migrations (requires the `migrate` CLI tool):

```bash
./scripts/migrate.sh up
```

#### 3. Build and Run Locally

```bash
go build -o urlshortener .
./urlshortener
```

The service should start on port `8080` by default.

#### 4. Or Run with Docker

```bash
docker build -t urlshortener .
docker run -p 8080:8080 urlshortener
```

## API Endpoints

- `POST /api/v1/shortly/url`  
  Shortens a URL.  
  **Body:** `{"long_url": "https://your-long-url"}`

- `GET /api/v1/shortly/urls`  
  List all shortened URLs.

- `GET /shortly/{id}`  
  Redirect to the original URL by short code.

- `GET /health`  
  Service health check.

- `GET /test-db`  
  Simple database operation tests (for development).

## Example Usage

Shorten a URL:

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"long_url":"https://example.com"}' \
  http://localhost:8080/api/v1/shortly/url
```

Redirect:

Visit `http://localhost:8080/shortly/{short_code}` in your browser.

## License

MIT
