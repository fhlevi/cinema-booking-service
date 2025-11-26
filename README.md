# Cinema Booking System

A microservice-based cinema ticket booking system built with Go, Gin framework, PostgreSQL, and Docker.

## Features

- **User Authentication**: Register and login with JWT tokens
- **Cinema Management**: 5 studios with 20 seats each
- **Online Booking**: Authenticated users can book tickets
- **Offline Booking**: Cashiers can create bookings for walk-in customers
- **QR Code Generation**: Each booking automatically generates a unique QR code
- **Seat Locking**: Reserved seats are immediately locked and unavailable to others
- **Ticket Validation**: QR codes can be validated at studio entrance

## Architecture

The system consists of 4 microservices:

1. **Auth Service** (Port 3001): Handles user authentication and JWT
2. **Cinema Service** (Port 3002): Manages studios and seats
3. **Booking Service** (Port 3003): Handles ticket booking and QR code generation
4. **API Gateway** (Port 3000): Routes requests and provides API documentation

## Getting Started

### Prerequisites
- Docker and Docker Compose

### Running the Application

1. Clone the repository:
```bash
git clone <repository-url>
cd cinema-booking
```

2. Start all services:
```bash
docker-compose up --build
```

3. Access the API:
- API Gateway: http://localhost:3000
- **Swagger UI Documentation**: http://localhost:3000/api/docs
- Swagger JSON: http://localhost:3000/docs/swagger.json

### Development Mode

To run individual services for development:

```bash
# Download dependencies for all services
go mod tidy

# Run a specific service
cd auth-service
go run main.go
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `POST /api/auth/verify` - Verify JWT token

### Cinema Management
- `GET /api/cinema/studios` - Get all studios
- `GET /api/cinema/studios/:id/seats` - Get studio seats

### Booking
- `POST /api/booking/online` - Create online booking (requires auth)
- `POST /api/booking/offline` - Create offline booking (cashier)
- `POST /api/booking/validate` - Validate QR code
- `GET /api/booking/my-bookings` - Get user bookings (requires auth)

## API Usage Examples

### 1. Register User
```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 3. Get Studio Data
```bash
curl http://localhost:3000/api/cinema/studios
```

### 4. Get Studio Seats
```bash
curl http://localhost:3000/api/cinema/studios/1/seats
```

### 5. Online Booking
```bash
curl -X POST http://localhost:3000/api/booking/online \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "studioId": 1,
    "seatIds": [1, 2, 3]
  }'
```

### 6. Offline Booking (Cashier)
```bash
curl -X POST http://localhost:3000/api/booking/offline \
  -H "Content-Type: application/json" \
  -d '{
    "studioId": 1,
    "seatIds": [4, 5],
    "customerName": "Jane Doe",
    "customerEmail": "jane@example.com"
  }'
```

### 7. Validate QR Code
```bash
curl -X POST http://localhost:3000/api/booking/validate \
  -H "Content-Type: application/json" \
  -d '{
    "bookingCode": "BOOKING_CODE_FROM_QR"
  }'
```

## Database Schema

### Users Table
- id (Primary Key)
- email (Unique)
- password (Hashed)
- name
- role (default: 'customer')
- created_at

### Studios Table
- id (Primary Key)
- name
- total_seats (default: 20)
- created_at

### Seats Table
- id (Primary Key)
- studio_id (Foreign Key)
- seat_number
- is_available (Boolean)

### Bookings Table
- id (Primary Key)
- booking_code (Unique)
- user_id (Nullable for offline bookings)
- user_name
- user_email
- studio_id
- seat_ids (Array)
- qr_code (Base64 encoded)
- booking_type ('online' or 'offline')
- status ('active' or 'used')
- created_at

## Testing

The project includes comprehensive unit tests for all services covering:
- **Handler validation** - Request/response validation and error handling
- **JWT utilities** - Token generation and validation
- **QR code generation** - QR code creation and data structure validation
- **API Gateway** - Routing, CORS, and proxy functionality
- **Service integration** - Cross-service communication patterns

### Running Tests

**Run all tests across all services:**
```bash
# Windows
run-all-tests.bat

# Linux/macOS
chmod +x run-all-tests.sh
./run-all-tests.sh
```

**Run tests for individual services:**
```bash
# Auth Service
cd auth-service
go test -v ./...

# Cinema Service
cd cinema-service
go test -v ./...

# Booking Service
cd booking-service
go test -v ./...

# API Gateway
cd api-gateway
go test -v ./...
```

**Generate test coverage reports:**
```bash
# Windows
test-coverage.bat

# Individual service coverage
cd auth-service
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Test Structure

Each service includes:
- `main_test.go` - Basic integration tests
- `handlers/*_test.go` - Handler-specific unit tests
- `utils/*_test.go` - Utility function tests

**Test Dependencies:**
- `github.com/stretchr/testify` - Enhanced assertions and test utilities
- `net/http/httptest` - HTTP testing utilities
- `gin.TestMode` - Gin framework test mode

## Deployment

### Production Build with Docker
```bash
docker-compose -f docker-compose.prod.yml up --build
```

### Environment Variables
- `DATABASE_URL`: PostgreSQL connection string (format: postgres://user:password@host:port/dbname?sslmode=disable)
- `JWT_SECRET`: Secret key for JWT tokens
- `AUTH_SERVICE_URL`: Auth service URL
- `CINEMA_SERVICE_URL`: Cinema service URL
- `BOOKING_SERVICE_URL`: Booking service URL
- `PORT`: Service port (default: 8080)

## Security Features

- JWT-based authentication
- Password hashing with bcrypt
- CORS protection
- Input validation
- SQL injection prevention with prepared statements
- Secure HTTP headers

## Business Logic

1. **Cinema Setup**: System starts with 5 studios, each having 20 seats (A1-A20)
2. **Seat Reservation**: When booking is created, seats are immediately locked
3. **QR Code**: Contains booking code, user info, studio, seats, and timestamp
4. **Validation**: QR codes can only be used once and mark booking as 'used'
5. **Concurrent Booking**: Database transactions prevent double-booking

## Monitoring

Each service has a health check endpoint:
- Auth Service: http://localhost:3001/health
- Cinema Service: http://localhost:3002/health
- Booking Service: http://localhost:3003/health
- API Gateway: http://localhost:3000/health

## Go-Specific Features

- **Gin Framework**: Fast HTTP web framework
- **Goroutines**: Concurrent request handling
- **Built-in Testing**: Go's native testing framework
- **Static Compilation**: Single binary deployment
- **Memory Efficiency**: Low memory footprint
- **Fast Startup**: Quick service initialization