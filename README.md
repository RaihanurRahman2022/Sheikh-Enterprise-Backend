# Sheikh Enterprise Backend

A robust backend service built with Go (Gin + GORM) for managing sales and inventory.

## Features

- JWT Authentication & Role-based Authorization
- Product Management
  - CRUD operations
  - Bulk Import
  - Filtering and Sorting
  - Excel Export
- Sales Management
  - CRUD operations
  - Filtering and Sorting
  - Excel Export
- User Management
  - Profile Management
  - Password Management
- Analytics Dashboard
  - Daily/Monthly/Yearly Sales
  - Product Statistics
  - Sales Charts

## Tech Stack

- Go 1.21+
- Gin Web Framework
- GORM
- PostgreSQL
- JWT for Authentication

## Project Structure

```
.
├── cmd/
│   └── api/                    # Application entrypoint
├── config/                     # Configuration files
├── internal/
│   ├── models/                 # Database models
│   ├── handlers/               # HTTP handlers
│   ├── middleware/             # HTTP middleware
│   ├── repository/             # Database repositories
│   ├── services/               # Business logic
│   └── utils/                  # Utility functions
├── pkg/
│   ├── database/              # Database connection
│   ├── logger/                # Logging package
│   └── validator/             # Validation package
└── docs/                      # Documentation
```

## Getting Started

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Run `go mod download` to install dependencies
4. Run `go run cmd/api/main.go` to start the server

## API Documentation

API documentation is available at `/swagger/index.html` when running in development mode.

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=sheikh_enterprise
JWT_SECRET=your_jwt_secret
PORT=8080
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request # Sheikh-Enterprise-Backend
