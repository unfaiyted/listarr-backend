# Listarr Backend
A RESTful API service built with Go and documented with Swagger/OpenAPI.

## Prerequisites

- Go 1.20 or higher
- Docker (optional)
- Make (optional)

## Project Structure
├── cmd/
│   └── api/          # Application entrypoint
├── internal/
│   ├── api/         # API handlers and routes
│   ├── config/      # Configuration management
│   ├── middleware/  # HTTP middleware
│   ├── models/      # Data models
│   └── service/     # Business logic
├── pkg/             # Public packages
├── docs/           # Swagger documentation
└── scripts/        # Build and deployment scripts
```

## Getting Started

### Local Development

1. Clone the repository
```bash
git clone https://github.com/yourusername/listarr-backend.git
cd listarr-backend
```

2. Install dependencies
```bash
go mod download
```

3. Run the application
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### Using Docker

```bash
# Build the Docker image
docker build -t listarr-backend .

# Run the container
docker run -p 8080:8080 listarr-backend
```

## API Documentation

The API documentation is available through Swagger UI when the application is running:

```
http://localhost:8080/swagger/index.html
```

### Key Endpoints

- `GET /api/v1/health` - Health check endpoint
- `GET /api/v1/docs` - API documentation (Swagger UI)

## Configuration

The application can be configured using environment variables or a configuration file. See `.env.example` for available options.

## Development

### Adding New Endpoints

1. Define the route in `internal/api/routes.go`
2. Create handler in `internal/api/handlers/`
3. Update Swagger documentation using comments
4. Generate new Swagger docs:
```bash
swag init -g cmd/api/main.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

