# Go HTTP Service Template

A production-ready HTTP service template with built-in observability, security, extensibility features, and AI integration capabilities.

## ✨ Features

- 🤖 AI Chat integration with multiple provider support
- 🔄 Graceful shutdown handling
- 📊 Prometheus metrics built-in
- 📝 Structured logging with Zap
- 🔒 Security headers and rate limiting
- 🔍 Request tracing with unique IDs
- ⚡ Automatic CPU core utilization
- 🎯 Middleware-based architecture
- 🚦 Environment-based configuration

## 🏗 Project Structure

```
.
├── internal/
│   ├── handlers/     # HTTP request handlers
│   │   ├── base.go   # Base handler with common response methods
│   │   ├── hello.go  # Example handler implementation
│   │   └── chat.go   # AI chat handler
│   ├── middleware/   # HTTP middleware components
│   └── ai/          # AI provider implementations
│       ├── provider.go    # Provider interface and registry
│       └── openrouter.go  # OpenRouter provider implementation
├── pkg/
│   └── server/       # Core server package
├── main.go          # Application entry point
└── go.mod           # Go module definition
```

## 🚀 Quick Start

1. Clone the repository:

   ```bash
   git clone https://github.com/rezashahnazar/golang-ai-webservcie.git
   cd golang-ai-webservcie
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Set up environment variables:

   ```bash
   cp .env.example .env
   # Edit .env and add your OpenRouter API key
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

## ⚙️ Configuration

Configuration is handled through environment variables:

| Variable           | Description                  | Default                   |
| ------------------ | ---------------------------- | ------------------------- |
| PORT               | Server port                  | 8080                      |
| READ_TIMEOUT       | Read timeout (seconds)       | 15                        |
| WRITE_TIMEOUT      | Write timeout (seconds)      | 15                        |
| IDLE_TIMEOUT       | Idle timeout (seconds)       | 60                        |
| RATE_LIMIT_RPS     | Rate limit (requests/second) | 100                       |
| OPENROUTER_API_KEY | OpenRouter API key           | Required                  |
| OPENROUTER_MODEL   | OpenRouter model name        | anthropic/claude-3-sonnet |

## 🔌 API Endpoints

### Standard Endpoints

- `GET /hello` - Example endpoint returning a greeting
- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics endpoint

### AI Chat Endpoint

- `POST /v1/chat` - Stream AI chat responses

#### Chat Request Format

```json
{
  "messages": [{ "role": "user", "content": "Hello, how are you?" }],
  "provider": "openrouter", // optional, defaults to "openrouter"
  "model": "anthropic/claude-3-sonnet", // optional
  "max_tokens": 1000, // optional
  "temperature": 0.7 // optional
}
```

The response is streamed as Server-Sent Events (SSE) containing the AI's response chunks.

## 📊 Metrics

Available Prometheus metrics:

- `http_requests_total` - Total HTTP requests (labels: method, endpoint, status)
- `http_request_duration_seconds` - Request duration histogram (labels: method, endpoint)

## 🛠 Adding New AI Providers

1. Create a new provider in `internal/ai/`:

```go
package ai

type MyProvider struct {
    // provider-specific fields
}

func NewMyProvider(options *ProviderOptions) (Provider, error) {
    // Initialize your provider
}

func (p *MyProvider) ChatStream(ctx context.Context, messages []Message, options *ProviderOptions) (io.ReadCloser, error) {
    // Implement chat streaming
}

func init() {
    RegisterProvider("myprovider", NewMyProvider)
}
```

2. Use the provider in requests:

```json
{
  "messages": [...],
  "provider": "myprovider"
}
```

## 🔒 Security Features

- Rate limiting per endpoint
- Security headers:
  - X-Content-Type-Options
  - X-Frame-Options
  - X-XSS-Protection
  - Strict-Transport-Security
- Request ID tracking
- Panic recovery middleware
- Environment-based configuration
- Secure API key handling

## 🧪 Testing

Run the tests:

```bash
go test ./...
```

## 📦 Docker Support

Build and run with Docker:

```bash
docker build -t go-http-service .
docker run -p 8080:8080 --env-file .env go-http-service
```

### Development with Hot Reloading

For development, the project includes hot reloading support using Air. This allows you to make changes to your code and see them reflected immediately without manually restarting the server.

1. Start the development server with hot reloading:

```bash
docker compose up -d --build
```

2. View logs while developing:

```bash
docker compose logs -f api
```

The development setup includes:

- Automatic rebuilding when files change
- Live reloading of the application
- Volume mounting for instant code syncing
- Air configuration for efficient Go rebuilds

Configuration files:

- `.air.toml` - Air configuration for hot reloading
- `Dockerfile.dev` - Development-specific Docker setup
- `docker-compose.yml` - Development container orchestration

## 👥 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 📧 Contact

Reza Shahnazar - reza.shahnazar@gmail.com

Project Link: [https://github.com/rezashahnazar/golang-ai-webservcie](https://github.com/rezashahnazar/golang-ai-webservcie)
