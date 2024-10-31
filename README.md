# Money Tracker Backend

This is the backend service for the Money Tracker application, built with Go and Gin framework.

## Project Structure

The project follows a clean architecture and is organized as follows:


### Directory Explanations

- `cmd/`: Contains the main application entry points.
    - `main.go`: The main entry point of the application.

- `internal/`: Houses the core application code.
    - `auth/`: Authentication-related middleware and utilities.
    - `controllers/`: HTTP request handlers.
    - `dto/`: Data Transfer Objects for API requests and responses.
    - `models/`: Data models representing the application's entities.
    - `repositories/`: Data access layer for interacting with the database.
    - `services/`: Business logic layer.
    - `utils/`: Utility functions and helpers.

- `migrations/`: Contains SQL migration scripts for database schema changes.

- `pkg/`: Shared packages that could potentially be used by other projects.
    - `database/`: Database connection and initialization utilities.

- `.env`: Environment variable configuration file.

- `go.mod` and `go.sum`: Go module definition and checksum files.

## Key Components

1. **Main Application (`cmd/main.go`)**: Initializes the application, sets up the database connection, runs migrations, and defines API routes.

2. **Controllers**: Handle incoming HTTP requests and manage the flow of data between the client and the application's business logic.

3. **Services**: Implement the core business logic of the application.

4. **Repositories**: Manage data storage and retrieval, abstracting the database operations from the rest of the application.

5. **Models**: Define the structure of the data used in the application.

6. **Database Migrations**: Allow for version-controlled changes to the database schema.

7. **Authentication Middleware**: Ensures that certain routes are protected and only accessible to authenticated users.

## Getting Started

1. Clone the repository
2. Set up your `.env` file with necessary configurations
3. Run `go mod download` to install dependencies
4. Use `go run cmd/main.go` to start the server

For more detailed information on setting up and running the project, please refer to the project documentation.