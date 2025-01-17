# Golang AI Web Service

A production-ready HTTP web service built with Go, featuring graceful shutdown, Docker support, and comprehensive health checking.

## Features

- 🚀 High-performance HTTP server utilizing all available CPU cores
- 🔒 Secure configuration with timeout settings
- 🐳 Docker containerization with resource limits
- 💪 Graceful shutdown handling
- 🏥 Health check endpoint for monitoring
- 📝 JSON response formatting

## Prerequisites

- Go 1.22 or higher
- Docker (optional)
- pnpm (for frontend dependencies)

## Getting Started

### Local Development

1. Clone the repository:

```bash
git clone https://github.com/rezashahnazar/golang-ai-webservice.git
cd golang-ai-webservice
```

2. Run the server:

```bash
go run main.go
```

The server will start on port 8080.

### Docker Deployment

1. Build the Docker image:

```bash
docker compose build
```

2. Run the container:

```bash
docker compose up -d
```

## API Endpoints

### Hello World

- **URL**: `/hello`
- **Method**: GET
- **Response**: `{"message": "Hello, World!"}`

### Health Check

- **URL**: `/health`
- **Method**: GET
- **Response**: 200 OK

## Configuration

The service includes several production-ready configurations:

- Server timeouts:

  - Read timeout: 15 seconds
  - Write timeout: 15 seconds
  - Idle timeout: 60 seconds

- Docker resource limits:
  - CPU: 1 core
  - Memory: 512MB
  - Log rotation: 10MB max file size, 3 files maximum

## Development

### Project Structure

```
.
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── main.go
└── README.md
```

## Security

- Docker container runs with no-new-privileges security option
- Proper timeout configurations to prevent DOS attacks
- Alpine-based minimal Docker image to reduce attack surface

## Contributing

Please feel free to submit issues, fork the repository, and create pull requests for any improvements.

## Contact

- Author: Reza Shahnazar
- GitHub: [@rezashahnazar](https://github.com/rezashahnazar)
- Email: reza.shahnazar@gmail.com

## License

This project is licensed under the MIT License - see the LICENSE file for details.
