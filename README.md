# Fledge API - Travel Booking Platform

## Overview
Fledge is a comprehensive travel booking platform that allows users to search and book flights, hotels, and vacation packages. This repository contains the backend API built with Go, following clean architecture principles.

## Features
- 🔐 User authentication and authorization
- ✈️ Flight search and booking
- 🏨 Hotel search and booking
- 🌴 Vacation package management
- 💳 Booking management
- 📱 User profile and preferences
- 📊 Booking history and analytics

## Technology Stack
- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Documentation**: Swagger/OpenAPI
- **Testing**: Go testing package with testify
- **CI/CD**: GitHub Actions

## Project Structure
```
.
├── cmd
│   └── api                 # Application entrypoint
├── internal
│   ├── config             # Configuration
│   ├── domain             # Business logic and entities
│   │   ├── entity
│   │   └── repository
│   ├── handler            # HTTP handlers
│   ├── middleware         # Middleware components
│   ├── service            # Business logic implementation
│   └── util               # Utility functions
├── pkg
│   ├── errors             # Custom error definitions
│   └── validator          # Custom validators
└── docs                   # Documentation
```

## Prerequisites
- Go 1.21 or higher
- PostgreSQL 13 or higher
- Docker (optional)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/yourusername/fledge-api.git
cd fledge-api
```

2. Create and configure the .env file:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Install dependencies:
```bash
go mod download
```

4. Run the database migrations:
```bash
make migrate
```

5. Start the server:
```bash
make run
```

## API Documentation

### Authentication Endpoints
- `POST /auth/signup` - Register a new user
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh access token

### Flight Endpoints
- `GET /api/flights/search` - Search available flights
- `GET /api/flights/{id}` - Get flight details
- `POST /api/flights/{id}/book` - Book a flight

### Hotel Endpoints
- `GET /api/hotels/search` - Search available hotels
- `GET /api/hotels/{id}` - Get hotel details
- `POST /api/hotels/{id}/book` - Book a hotel

### Vacation Package Endpoints
- `GET /api/packages` - List vacation packages
- `GET /api/packages/{id}` - Get package details
- `POST /api/packages/{id}/book` - Book a package

### Booking Management
- `GET /api/bookings` - List user bookings
- `GET /api/bookings/{id}` - Get booking details
- `PATCH /api/bookings/{id}` - Update booking
- `DELETE /api/bookings/{id}` - Cancel booking

### User Profile
- `GET /api/profile` - Get user profile
- `PUT /api/profile` - Update user profile
- `GET /api/profile/preferences` - Get travel preferences
- `PUT /api/profile/preferences` - Update preferences

## Development

### Running Tests
```bash
make test
```

### Running Linter
```bash
make lint
```

### Building for Production
```bash
make build
```

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.