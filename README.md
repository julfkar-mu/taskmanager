# Task Manager API

A RESTful API for managing tasks built with Go and Gin framework. This project demonstrates clean architecture principles, comprehensive testing, and Go best practices.

## Features

- ✅ CRUD operations for tasks
- ✅ Input validation and error handling
- ✅ Clean architecture with separation of concerns
- ✅ Comprehensive unit tests with high coverage
- ✅ RESTful API design
- ✅ Concurrent-safe in-memory storage
- ✅ Docker support
- ✅ CI/CD with GitHub Actions
- ✅ API documentation with Swagger annotations

## Architecture

The project follows clean architecture principles with the following layers:

```
├── controllers/     # HTTP handlers and request/response handling
├── services/        # Business logic layer
├── repository/      # Data access layer
├── models/          # Domain models and entities
├── errors/          # Custom error types
├── constants/       # Application constants
└── testutils/       # Test utilities and helpers
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tasks` | Get all tasks |
| GET | `/api/v1/tasks/{id}` | Get task by ID |
| POST | `/api/v1/tasks` | Create a new task |
| PUT | `/api/v1/tasks/{id}` | Update a task |
| DELETE | `/api/v1/tasks/{id}` | Delete a task |
| GET | `/health` | Health check |

## Task Model

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Complete project documentation",
  "description": "Write comprehensive documentation for the API",
  "status": "Pending",
  "priority": "High",
  "dueDate": "2024-12-31T23:59:59Z",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
  "assignedTo": "john.doe@example.com"
}
```

### Task Status Values
- `Pending` - Task is not started
- `InProgress` - Task is currently being worked on
- `Completed` - Task is finished
- `Cancelled` - Task was cancelled

### Task Priority Values
- `Low` - Low priority task
- `Medium` - Medium priority task
- `High` - High priority task

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/taskmanager.git
cd taskmanager
```

2. Download dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

### Using Docker

1. Build the Docker image:
```bash
docker build -t taskmanager .
```

2. Run the container:
```bash
docker run -p 8080:8080 taskmanager
```

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# Run benchmarks
make bench
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run all checks
make check
```

### Available Make Commands

- `make build` - Build the application
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make lint` - Run linter
- `make fmt` - Format code
- `make clean` - Clean build artifacts
- `make run` - Run the application
- `make docker-build` - Build Docker image
- `make docker-run` - Run Docker container

## Testing

The project includes comprehensive unit tests for all layers:

- **Models**: Validation logic and business rules
- **Repository**: Data access layer with concurrency tests
- **Services**: Business logic with mocked dependencies
- **Controllers**: HTTP handlers with mocked services

Test coverage is maintained above 90% and includes:
- Unit tests for all functions
- Integration tests for API endpoints
- Concurrency tests for thread safety
- Error handling tests
- Validation tests

## Error Handling

The API uses custom error types for better error handling:

- `ValidationError` - Input validation errors (400 Bad Request)
- `AppError` - Application errors with HTTP status codes
- `NotFoundError` - Resource not found errors (404 Not Found)

## API Examples

### Create a Task

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project documentation",
    "description": "Write comprehensive documentation for the API",
    "status": "Pending",
    "priority": "High",
    "assignedTo": "john.doe@example.com"
  }'
```

### Get All Tasks

```bash
curl http://localhost:8080/api/v1/tasks
```

### Update a Task

```bash
curl -X PUT http://localhost:8080/api/v1/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated task title",
    "status": "InProgress"
  }'
```

### Delete a Task

```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/{id}
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gin](https://gin-gonic.com/) - HTTP web framework
- [Testify](https://github.com/stretchr/testify) - Testing toolkit
- [UUID](https://github.com/google/uuid) - UUID generation

