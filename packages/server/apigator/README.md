# Apigator Service

Apigator is a Go-based authentication service that validates tokens using PostgreSQL and Neo4j databases. It is designed to provide a robust and scalable solution for managing authentication in distributed systems.

## Features

- **Token Validation**: Validates tokens using repositories from PostgreSQL and Neo4j.
- **Tracing**: Integrated with Jaeger for distributed tracing.
- **Logging**: Utilizes Logrus for structured logging.
- **Configuration**: Supports environment-based configuration using `godotenv`.

## Getting Started

### Prerequisites

- Go 1.20 or later
- PostgreSQL and Neo4j databases

### Running the Application

1. Ensure that your environment variables are set up correctly. You can use a `.env` file for configuration.
2. Run the application with the following command:

   ```bash
   go run main.go
   ```

3. The service will start and listen on the port specified in your configuration.

### API Endpoints

- **GET /validate**: Validates a token using the configured repositories.

## Project Structure

- `main.go`: The main entry point for the application.
- `go.mod`: The module file for dependency management.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
