# Parking Places Booking System

A comprehensive microservices-based parking places booking platform built with Go, featuring distributed tracing, monitoring, event-driven notifications, and a modern React frontend.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Microservices](#microservices)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Testing](#testing)
- [Monitoring and Observability](#monitoring-and-observability)
- [HTTPS Setup](#https-setup)
- [Project Structure](#project-structure)
- [Troubleshooting](#troubleshooting)
- [Security Considerations](#security-considerations)
- [Production Deployment](#production-deployment)

## Overview

This project implements a parking places booking system using a microservices architecture. The system supports two user roles: drivers who book parking spaces and parking owners who manage their parking facilities. The platform provides REST APIs, gRPC services, a modern React frontend, and a Telegram bot interface.

### Key Features

- Clean layered architecture with Handler, Service, and Repository layers
- Domain-driven design with separate domain models
- Interface-based design for testability
- Standardized error handling across services
- Multi-tenant architecture with role-based access control
- Distributed tracing with Jaeger
- Metrics collection with Prometheus
- Event-driven notifications via Kafka
- OAuth2/OIDC authentication with Keycloak
- REST APIs for external communication and gRPC for inter-service calls
- Full Docker Compose orchestration
- Modern React frontend with multi-language support
- Telegram bot interface for user interaction
- Secure payment processing with promocode system
- Support for multiple parking types: outdoor, covered, underground, and multi-level

## Architecture

The system follows a microservices architecture with a clean layered design:

- **Handler Layer**: HTTP request/response handling and domain-to-API mapping
- **Service Layer**: Business logic, validation, and orchestration
- **Repository Layer**: Data access with interface-based design
- **Domain Models**: Core business entities separate from API models

Each service follows the same consistent structure:

```
Handler → Service → Repository → Database
```

## Technology Stack

### Core Technologies

- **Language**: Go 1.23.1
- **API Framework**:
  - REST: go-swagger (OpenAPI 2.0)
  - gRPC: Protocol Buffers v3
- **Database**: PostgreSQL 16.4
- **Authentication**: Keycloak (OAuth2/OIDC)
- **Message Queue**: Apache Kafka 7.3.0 with Zookeeper
- **Containerization**: Docker and Docker Compose

### Frontend Technologies

- **Framework**: React 18.3 with Vite 5
- **Styling**: Tailwind CSS 3
- **HTTP Client**: Axios
- **Routing**: React Router 6
- **Internationalization**: i18next

### Observability and Monitoring

- **Distributed Tracing**: Jaeger
- **Metrics**: Prometheus with Grafana dashboards
- **Logging**: Structured logging across all services

### Libraries and Frameworks

- **Database Driver**: jackc/pgx/v5
- **Keycloak Client**: Nerzal/gocloak/v13
- **Kafka Client**: segmentio/kafka-go
- **Telegram Bot**: go-telegram-bot-api/v5
- **OpenTelemetry**: OTEL SDK for tracing
- **Configuration**: Centralized config management with validation

## Microservices

### 1. Auth Service (Port 8800)

Responsibility: User authentication and authorization

Features:
- User registration with role assignment (driver/owner)
- Login with JWT token generation
- Password change functionality
- Keycloak integration for identity management
- Token validation for protected endpoints
- Admin user management

API Endpoints:
- `POST /login` - User authentication
- `POST /register` - New user registration
- `POST /change-password` - Password modification
- `GET /user` - Get user information (admin only)
- `GET /metrics` - Prometheus metrics

Database: Uses Keycloak's database for user management

### 2. Parking Service (REST: Port 8888, gRPC: Port 50051)

Responsibility: Parking place management and information retrieval

Features:
- CRUD operations for parking places
- Search parking by city, name, or type
- Role-based access control (owners manage their parking places)
- Dual API exposure (REST and gRPC)
- gRPC service for internal service-to-service communication
- Hourly rate-based pricing model
- Domain models with validation

API Endpoints:
- `GET /parking` - Search parking places with filters
- `POST /parking` - Create new parking place (owners only)
- `GET /parking/{parking_id}` - Get parking place details
- `PUT /parking/{parking_id}` - Update parking place (owner only)
- `DELETE /parking/{parking_id}` - Delete parking place (owner only)
- `GET /metrics` - Prometheus metrics

gRPC Service:
- `GetParkingPlace(ParkingPlaceRequest)` - Retrieve parking place information

Database: `parking_db`

Schema:
```sql
parking_places (id, name, city, address, parking_type, hourly_rate, capacity, owner_id)
```

### 3. Booking Service (Port 8880)

Responsibility: Booking management and lifecycle

Features:
- Create bookings with date validation
- Automatic payment processing on booking creation
- Booking status management (Waiting, Confirmed, Canceled)
- Retrieve bookings by ID or parking place
- Calculate total cost based on hourly rate and duration
- gRPC client to fetch parking place information
- gRPC client for payment processing
- Role-based access (drivers book, owners manage)
- Automatic refunds on booking cancellation

API Endpoints:
- `POST /booking` - Create new booking (drivers)
- `GET /booking` - Get bookings by parking place (owners)
- `GET /booking/{booking_id}` - Get booking details
- `PUT /booking/{booking_id}` - Update booking status
- `DELETE /booking/{booking_id}` - Cancel booking with refund
- `GET /metrics` - Prometheus metrics

Database: `booking_db`

Schema:
```sql
bookings (id, date_from, date_to, parking_place_id, full_cost, status, user_id)
```

### 4. Payment Service (REST: Port 8890, gRPC: Port 50052)

Responsibility: Financial operations and billing

Features:
- User balance management
- Transaction processing (charge drivers, pay owners)
- Refund processing
- Promocode system:
  - Activate promocodes to add balance
  - Generate promocodes from user balance (withdrawal)
  - Admin creation of custom promocodes
- Transaction history
- Atomic transactions with database locking
- Overflow protection for balance operations
- gRPC service for internal payment processing

API Endpoints:
- `GET /payment/balance` - Get user balance
- `GET /payment/transactions` - Get transaction history
- `POST /payment/promocode/activate` - Activate a promocode
- `POST /payment/promocode/generate` - Generate promocode from balance
- `POST /payment/promocode` - Create promocode (admin only)
- `GET /payment/promocode/{code}` - Get promocode information
- `GET /metrics` - Prometheus metrics

gRPC Service:
- `ProcessTransaction(TransactionRequest)` - Process payment transaction
- `ProcessRefund(RefundRequest)` - Process refund transaction

Database: `payment_db`

Schema:
```sql
balances (user_id, balance, currency)
transactions (id, user_id, amount, type, status, booking_id, created_at)
promocodes (code, amount, usage_limit, used_count, expires_at, created_by)
```

### 5. Notification Service

Responsibility: Asynchronous notification handling

Features:
- Kafka consumer for notification events
- Telegram notification integration
- Graceful error handling and service resilience
- Extensible handler architecture
- Continues processing messages even when individual notifications fail

Kafka Integration:
- Topic: Configurable via environment
- Group ID: notification-service
- Message handlers for different notification types

### 6. Telegram Bot

Responsibility: Conversational interface for the platform

Features:
- Simple command-based interface
- User authentication through the bot
- Balance viewing
- Booking management
- Parking place management
- Role-specific menu options
- Session management per user

Supported Commands:
- `/start` - Initialize bot and display Telegram ID
- `/help` - Show available commands
- `/login` - Authenticate user (login and password)
- `/balance` - View account balance
- `/bookings` - View current bookings (drivers)
- `/parkings` - View owned parking places (owners)

### 7. Frontend (Port 3000) - Standalone Deployment

Responsibility: Modern web interface for the platform

Features:
- React.js with Vite for fast development
- Tailwind CSS for beautiful, responsive UI
- Multi-language support (English/Russian)
- Role-based dashboards (Driver and Owner)
- Real-time data updates
- Mobile-responsive design
- JWT token authentication
- Fully decoupled from backend
- Admin panel with monitoring tools

Driver Features:
- Search parking places with filters
- Create and manage bookings
- View booking history and status
- Balance management
- Promocode activation and generation

Owner Features:
- Create and manage parking places
- View all bookings for owned parkings
- Access admin monitoring tools

Common Features:
- Admin panel with links to Jaeger, Prometheus, and Grafana
- Monitoring and service health dashboards
- Promocode management (admin)

Deployment:
- Standalone Docker container
- Communicates with backend via HTTPS (nginx reverse proxy on host)
- Can be deployed on separate server
- Environment variable: `VITE_API_BASE_URL`

## Prerequisites

Before running the project, ensure you have the following installed:

- Docker: Version 20.10 or higher
- Docker Compose: Version 2.0 or higher
- Make (optional): For using Makefile commands
- Go 1.23.1+ (for local development)
- Protocol Buffers Compiler (for gRPC code generation)
- go-swagger (for API code generation)
- Python 3 (for running integration tests)

## Quick Start

### Initial Setup

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd parking_net
   ```

2. Configure environment variables:
   ```bash
   cp .env-example .env
   ```

3. Edit `.env` file with your configuration. Important values to set:
   - `KEYCLOAK_CLIENT_SECRET`: Set a secure random string (e.g., use `openssl rand -hex 32`)
   - `KEYCLOAK_ADMIN_PASSWORD`: Change from default if needed
   - `POSTGRES_PASSWORD`: Change from default if needed
   - `TELEGRAM_API_KEY`: Set your Telegram bot token if you need Telegram notifications
   - `INTERNAL_SERVICE_TOKEN`: Set a secure token for inter-service gRPC communication

   Security Note: The `.env` file is in `.gitignore` and will not be committed to Git. Never commit sensitive credentials.

4. Start all services:
   ```bash
   docker-compose up -d
   ```

   This will:
   - Start all backend services
   - Start the frontend on port 3000
   - Automatically run setup service after database and Keycloak are ready
   - Configure Keycloak client secret from your `.env` file
   - Create all required database tables
   - Set up Kafka configuration

   The setup runs automatically as part of docker-compose. No manual steps needed.

5. Verify setup:
   ```bash
   docker-compose logs setup
   docker-compose ps
   ```

   All services should be running. The setup service will show as "exited" after successful completion, which is normal.

6. Access the application:
   - Frontend: http://localhost:3000
   - Backend APIs: http://localhost:8800 (auth), http://localhost:8888 (parking), etc.
   - Keycloak Admin: http://localhost:8080
   - Jaeger UI: http://localhost:16686
   - Prometheus: http://localhost:9090
   - Grafana: http://localhost:3000 (if configured)

7. Run backend integration tests:
   ```bash
   python3 tests/integration_test.py
   ```

   Or using Make:
   ```bash
   make test
   ```

8. Run frontend tests:
   ```bash
   cd frontend
   npm install
   npm test
   ```

### Frontend Setup (Standalone)

The frontend can be deployed on a separate server from backend:

**Example: Backend on 192.168.1.100, Frontend on 192.168.1.101**

1. On backend server (192.168.1.100):
   ```bash
   docker-compose up -d
   ```
   Services are available on their direct ports (8800, 8888, 8880, 8890)
   For production, configure nginx reverse proxy on host for HTTPS (see [HTTPS Setup](#https-setup))

2. On frontend server (192.168.1.101):
   ```bash
   cd frontend
   docker build --build-arg VITE_API_BASE_URL=http://192.168.1.100 -t parking-frontend .
   docker run -d -p 3000:80 parking-frontend
   ```
   Frontend will be available at http://192.168.1.101:3000

3. For development:
   ```bash
   cd frontend
   echo "VITE_API_BASE_URL=http://192.168.1.100" > .env
   npm install
   npm run dev
   ```

### How Automated Setup Works

The setup is fully automated through a `setup` service in docker-compose:

1. Setup Service: Automatically runs after `db` and `keycloak` services are healthy
2. No Manual Steps: Everything happens automatically when you run `docker-compose up -d`
3. Works Everywhere: No Makefile needed - works in any environment with Docker

To check setup progress:
```bash
docker-compose logs -f setup
```

### Available Make Commands (Optional)

Make commands are optional convenience wrappers. Everything works with just `docker-compose`:

```bash
make help      # Show all available commands
make setup     # Initial setup (copy .env-example to .env)
make up        # Start all services (setup runs automatically)
make down      # Stop all services
make restart   # Restart all services
make test      # Run integration tests
make logs      # Show logs from all services
make ps        # Show status of all services
make clean     # Stop services and remove volumes
```

Note: Make is optional. You can use `docker-compose` commands directly - setup runs automatically.

## Configuration

### Environment Variables

Key environment variables (see `.env-example` for complete list):

**PostgreSQL:**
- `POSTGRES_USER`: Database user (default: postgres)
- `POSTGRES_PASSWORD`: Database password
- `POSTGRES_PORT`: Database port (default: 5432)
- `PARKING_DB_NAME`: Parking database name (default: parking_db)
- `BOOKING_DB_NAME`: Booking database name (default: booking_db)
- `PAYMENT_DB_NAME`: Payment database name (default: payment_db)
- `AUTH_DB_NAME`: Auth database name (default: auth_db)
- `TELEGRAM_DB_NAME`: Telegram database name (default: telegram_db)

**Service Ports:**
- `PARKING_REST_PORT`: Parking REST API port (default: 8888)
- `PARKING_GRPC_PORT`: Parking gRPC port (default: 50051)
- `BOOKING_REST_PORT`: Booking API port (default: 8880)
- `PAYMENT_REST_PORT`: Payment API port (default: 8890)
- `PAYMENT_GRPC_PORT`: Payment gRPC port (default: 50052)
- `AUTH_REST_PORT`: Auth API port (default: 8800)

**Keycloak:**
- `KEYCLOAK_PORT`: Keycloak port (default: 8080)
- `KEYCLOAK_CLIENT`: Client name (default: parking-auth)
- `KEYCLOAK_REALM`: Realm name (default: parking-users)
- `KEYCLOAK_CLIENT_SECRET`: Client secret (must be set)
- `KEYCLOAK_ADMIN`: Admin username (default: admin)
- `KEYCLOAK_ADMIN_PASSWORD`: Admin password
- `KEYCLOAK_FRONTEND_URL`: Keycloak frontend URL (for HTTPS)

**Kafka:**
- `KAFKA_BROKER`: Kafka broker address (default: kafka:9092)
- `KAFKA_TOPIC`: Kafka topic name
- `KAFKA_GROUP_ID`: Consumer group ID

**Telegram:**
- `TELEGRAM_API_KEY`: Telegram bot token

**Internal Service Authentication:**
- `INTERNAL_SERVICE_TOKEN`: Token for inter-service gRPC communication

**Domain Configuration (for HTTPS setup):**
- `FRONTEND_DOMAIN`: Frontend domain (e.g., parking-net.space)
- `BACKEND_DOMAIN`: Backend domain (e.g., backend.parking-net.space)
- `JAEGER_SUBDOMAIN`: Jaeger subdomain
- `GRAFANA_SUBDOMAIN`: Grafana subdomain
- `KEYCLOAK_SUBDOMAIN`: Keycloak subdomain
- `PROMETHEUS_SUBDOMAIN`: Prometheus subdomain
- `FRONTEND_IP`: Frontend server IP
- `BACKEND_IP`: Backend server IP

### Database Initialization

The system automatically creates five separate databases on startup:
- `parking_db` - Parking service data
- `booking_db` - Booking service data
- `payment_db` - Payment service data
- `auth_db` - Keycloak authentication data
- `telegram_db` - Telegram bot user data

Database schemas are initialized via SQL scripts in `scripts/init_sql/`:
- `init_parking.sql` - Parking places table
- `init_booking.sql` - Bookings table
- `init_payment.sql` - Balances, transactions, and promocodes tables
- `init_telegram.sql` - Telegram bot user data

### Keycloak Setup

Keycloak is pre-configured with a realm export (`keycloak/config/realm-export.json`). The configuration includes:
- Realm: `parking-users`
- Client: `parking-auth`
- Roles: `driver`, `owner`, `admin`

The setup service automatically configures the client secret from your `.env` file.

### Service Communication

- **External → Services**: REST APIs (Swagger/OpenAPI)
- **Booking → Parking**: gRPC (for parking place information retrieval)
- **Booking → Payment**: gRPC (for payment processing)
- **Services → Keycloak**: REST API for token validation
- **Services → Jaeger**: OTLP for trace export
- **Services → Kafka**: For event publishing/consumption
- **Inter-service gRPC**: Authenticated with `INTERNAL_SERVICE_TOKEN`

## API Documentation

API specifications are defined using OpenAPI 2.0 (Swagger):

- **Auth API Specification**: `auth/api/swagger/auth.yaml`
- **Parking API Specification**: `parking/api/swagger/parking.yaml`
- **Booking API Specification**: `booking/api/swagger/booking.yaml`
- **Payment API Specification**: `payment/api/swagger/payment.yaml`

### Authentication

Most endpoints require authentication. Include the token in the `api_key` header:

```bash
curl -H "api_key: YOUR_JWT_TOKEN" http://localhost:8888/parking
```

### Example API Usage

#### Register a New User

```bash
curl -X POST http://localhost:8800/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "login": "testuser",
    "password": "Password123",
    "role": "driver",
    "telegram_id": 123456789
  }'
```

#### Login

```bash
curl -X POST http://localhost:8800/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "testuser",
    "password": "Password123"
  }'
```

#### Create a Parking Place (Owner)

```bash
curl -X POST http://localhost:8888/parking \
  -H "Content-Type: application/json" \
  -H "api_key: YOUR_TOKEN" \
  -d '{
    "name": "Central Parking",
    "city": "Moscow",
    "address": "Red Square 1",
    "parking_type": "underground",
    "hourly_rate": 150,
    "capacity": 200
  }'
```

#### Search Parking Places

```bash
curl "http://localhost:8888/parking?city=Moscow&parking_type=underground"
```

#### Create a Booking (Driver)

```bash
curl -X POST http://localhost:8880/booking \
  -H "Content-Type: application/json" \
  -H "api_key: YOUR_TOKEN" \
  -d '{
    "parking_place_id": 1,
    "date_from": "2024-12-01T10:00:00Z",
    "date_to": "2024-12-05T18:00:00Z"
  }'
```

#### Activate Promocode

```bash
curl -X POST http://localhost:8890/payment/promocode/activate \
  -H "Content-Type: application/json" \
  -H "api_key: YOUR_TOKEN" \
  -d '{
    "code": "PROMO123"
  }'
```

#### Get Balance

```bash
curl -H "api_key: YOUR_TOKEN" http://localhost:8890/payment/balance
```

## Development

### Code Generation

The project uses code generation from OpenAPI specs and Protocol Buffers:

#### Generate All Code

```bash
make codegen
```

#### Generate from Swagger/OpenAPI

```bash
make swagger_generate
# Or manually:
./scripts/generate_from_swagger.sh
```

This generates:
- Server boilerplate
- Request/response models
- Validation logic
- API handlers structure

#### Generate gRPC Code

```bash
make grpc_generate
```

This generates Go code from `.proto` files for:
- Parking service gRPC server
- Payment service gRPC server
- Booking service gRPC clients

### Local Development Setup

1. Install dependencies:
   ```bash
   # For each service
   cd parking && go mod download
   cd ../booking && go mod download
   cd ../auth && go mod download
   cd ../payment && go mod download
   ```

2. Run services locally:
   ```bash
   # Set environment variables
   source .env

   # Run a service
   cd parking && go run cmd/main.go
   ```

3. Run with hot reload (using air or similar):
   ```bash
   air -c .air.toml
   ```

## Testing

### Integration Tests

The project includes comprehensive integration tests that verify the complete system flow:

```bash
python3 tests/integration_test.py
```

Or using Make:
```bash
make test
```

The integration tests cover:
- User registration and authentication
- Parking CRUD operations
- Parking search and filtering
- Booking creation and management
- Payment processing and promocodes
- Authorization checks
- Error handling (404, 403, 401)

### Unit Tests

Run unit tests for a specific service:

```bash
# Test a specific service
cd parking && go test ./...

# Test with coverage
go test -cover ./...

# Test with race detection
go test -race ./...
```

### Frontend Tests

Run frontend unit tests:

```bash
cd frontend
npm install
npm test
```

Run frontend E2E tests:

```bash
cd frontend
npm run test:e2e
```

### Load Testing

Use tools like `hey`, `wrk`, or `k6`:

```bash
# Example with hey
hey -n 10000 -c 100 http://localhost:8888/parking
```

## Monitoring and Observability

### Prometheus Metrics

All services expose metrics at `/metrics` endpoint:

- Request count: HTTP request totals
- Request duration: Latency histograms
- Error count: Failed request totals
- Custom business metrics: Service-specific metrics

Access Prometheus: http://localhost:9090

Example Queries:
```promql
# Request rate for parking service
rate(http_requests_total{job="parking"}[5m])

# 95th percentile latency
histogram_quantile(0.95, http_request_duration_seconds_bucket)

# Error rate
rate(http_requests_errors_total[5m])
```

### Distributed Tracing with Jaeger

All services are instrumented with OpenTelemetry.

Access Jaeger UI: http://localhost:16686

Trace Features:
- Request flow across services
- Service dependencies
- Latency breakdown
- Error tracking

### Grafana Dashboards

Pre-configured Grafana dashboards are available for:
- Service overview
- Individual service metrics (auth, parking, booking)
- System health monitoring

Access Grafana: http://localhost:3000 (default credentials: admin/admin)

### Logging

Services use structured logging with contextual information:
- Request ID
- User ID
- Service name
- Trace ID (for correlation with Jaeger)

View logs:
```bash
docker compose logs -f <service-name>
```

## HTTPS Setup

This guide covers setting up HTTPS for both frontend and backend services using Let's Encrypt SSL certificates.

### Prerequisites

- Domain names pointing to your server IPs
- Docker and Docker Compose installed
- Ports 80 and 443 open on your servers
- Root or sudo access to the servers

### Configuration

Before setup, configure domains in `.env`:

```bash
# Domain Configuration
FRONTEND_DOMAIN=parking-net.space
BACKEND_DOMAIN=backend.parking-net.space
JAEGER_SUBDOMAIN=jaeger.backend.parking-net.space
GRAFANA_SUBDOMAIN=grafana.backend.parking-net.space
KEYCLOAK_SUBDOMAIN=keycloak.backend.parking-net.space
PROMETHEUS_SUBDOMAIN=prometheus.backend.parking-net.space
KEYCLOAK_FRONTEND_URL=https://keycloak.backend.parking-net.space

# Server IPs (for DNS configuration)
FRONTEND_IP=158.160.159.53
BACKEND_IP=158.160.131.173
```

**Note:** Replace these values with your actual domains and IPs.

### Frontend HTTPS Setup

1. Install Nginx and Certbot:
   ```bash
   sudo apt update
   sudo apt install nginx certbot python3-certbot-nginx -y
   ```

2. Configure Nginx:
   ```bash
   sudo cp frontend/nginx-host.conf.example /etc/nginx/sites-available/parking-frontend
   sudo ln -s /etc/nginx/sites-available/parking-frontend /etc/nginx/sites-enabled/
   sudo mkdir -p /var/www/certbot
   ```
   
   Update `server_name` in config to match your `FRONTEND_DOMAIN`.

3. Test and Reload Nginx:
   ```bash
   sudo nginx -t
   sudo systemctl reload nginx
   ```

4. Get SSL Certificate:
   ```bash
   sudo certbot --nginx -d ${FRONTEND_DOMAIN} -d www.${FRONTEND_DOMAIN}
   ```
   
   Follow prompts and choose to redirect HTTP to HTTPS.

5. Enable Auto-Renewal:
   ```bash
   sudo systemctl enable certbot.timer
   sudo systemctl start certbot.timer
   ```

### Backend HTTPS Setup

1. Set Up DNS:
   
   Add A records:
   - `${BACKEND_DOMAIN}` → `${BACKEND_IP}`
   - `${JAEGER_SUBDOMAIN}` → `${BACKEND_IP}`
   - `${GRAFANA_SUBDOMAIN}` → `${BACKEND_IP}`
   - `${KEYCLOAK_SUBDOMAIN}` → `${BACKEND_IP}`
   - `${PROMETHEUS_SUBDOMAIN}` → `${BACKEND_IP}`

2. Configure Nginx for Backend:
   ```bash
   sudo cp nginx/nginx-backend-https.conf.example /etc/nginx/sites-available/parking-backend
   sudo ln -s /etc/nginx/sites-available/parking-backend /etc/nginx/sites-enabled/
   ```
   
   Update config file:
   - Replace placeholders with actual values from `.env`
   - Update proxy_pass ports to match your service ports

3. Test and Reload Nginx:
   ```bash
   sudo nginx -t
   sudo systemctl reload nginx
   ```

4. Get SSL Certificates:
   ```bash
   sudo certbot --nginx -d ${BACKEND_DOMAIN} -d ${JAEGER_SUBDOMAIN} -d ${GRAFANA_SUBDOMAIN} -d ${KEYCLOAK_SUBDOMAIN} -d ${PROMETHEUS_SUBDOMAIN}
   ```

5. Update Frontend Configuration:
   
   Frontend automatically uses HTTPS URLs in production. Verify in `.env`:
   ```bash
   VITE_API_BASE_URL=https://${BACKEND_DOMAIN}
   VITE_BASE_HOST=${BACKEND_DOMAIN}
   ```

6. Rebuild Frontend:
   ```bash
   cd frontend
   docker-compose build
   docker-compose up -d
   ```

7. Update Keycloak Configuration:
   
   Set in `.env`:
   ```bash
   KEYCLOAK_FRONTEND_URL=https://${KEYCLOAK_SUBDOMAIN}
   ```
   
   Then restart keycloak:
   ```bash
   docker-compose restart keycloak
   ```

### Verification

Frontend:
```bash
curl -I https://${FRONTEND_DOMAIN}
```

Backend:
```bash
curl -I https://${BACKEND_DOMAIN}/auth/metrics
curl -I https://${JAEGER_SUBDOMAIN}
curl -I https://${GRAFANA_SUBDOMAIN}
curl -I https://${KEYCLOAK_SUBDOMAIN}
curl -I https://${PROMETHEUS_SUBDOMAIN}
```

All should return 200 OK.

### Troubleshooting HTTPS

**502 Bad Gateway:**
- Check services are running: `docker-compose ps`
- Check ports match in nginx config
- Check nginx error log: `sudo tail -f /var/log/nginx/error.log`

**Certificate Errors:**
- Wait 5 minutes after DNS changes
- Check DNS: `dig ${BACKEND_DOMAIN}`
- Verify certificates: `sudo certbot certificates`

**Mixed Content Errors:**
- Clear browser cache
- Verify `API_BASE_URL` uses HTTPS
- Check CSP headers allow backend domain

**Auto-Renewal:**
Certificates expire every 90 days. Auto-renewal is enabled by default:
```bash
sudo systemctl status certbot.timer
```

To manually renew:
```bash
sudo certbot renew
```

## Project Structure

```
parking_net/
├── api/proto/                  # Protocol Buffer definitions
├── auth/                       # Auth microservice
│   ├── cmd/                    # Service entry point
│   ├── internal/
│   │   ├── impl/              # Business logic implementations
│   │   ├── models/            # API models (generated)
│   │   └── restapi/           # Generated API code
│   └── api/swagger/           # OpenAPI specification
├── frontend/                   # React.js web interface
│   ├── src/
│   │   ├── components/        # Reusable React components
│   │   ├── context/           # React context providers
│   │   ├── pages/             # Page components
│   │   ├── services/          # API service layer
│   │   └── config/            # Configuration files
│   ├── public/                # Static assets
│   ├── Dockerfile             # Frontend Docker config
│   └── package.json           # NPM dependencies
├── booking/                    # Booking microservice
│   ├── cmd/
│   ├── internal/
│   │   ├── repository/        # Data access layer
│   │   ├── database_service/  # Database operations
│   │   ├── grpc/              # gRPC client & generated code
│   │   ├── models/            # API models (generated)
│   │   └── restapi/           # Generated API code
│   └── api/swagger/
├── parking/                    # Parking microservice
│   ├── cmd/
│   │   ├── grpc/              # gRPC server startup
│   │   └── rest/              # REST server startup
│   ├── internal/
│   │   ├── repository/        # Data access layer
│   │   ├── service/           # Business logic layer
│   │   ├── handlers/          # HTTP handlers
│   │   ├── di/                # Dependency injection container
│   │   ├── grpc/              # gRPC server & generated code
│   │   ├── models/            # API models (generated)
│   │   └── restapi/           # Generated API code
│   └── api/swagger/
├── payment/                    # Payment microservice
│   ├── cmd/
│   │   ├── grpc/              # gRPC server startup
│   │   └── rest/              # REST server startup
│   ├── internal/
│   │   ├── database_service/  # Database operations
│   │   ├── grpc/              # gRPC server & generated code
│   │   ├── models/            # API models (generated)
│   │   └── restapi/           # Generated API code
│   └── api/swagger/
├── notification/               # Notification microservice
│   ├── cmd/
│   └── internal/
│       ├── handlers/          # Message handlers
│       ├── server/            # Kafka consumer
│       └── services/          # Notification logic
├── telegram_bot/               # Telegram bot interface
│   ├── cmd/
│   ├── api_service/           # API client for services
│   ├── database_service/      # Bot's database operations
│   └── data_representation/   # Message formatting
├── pkg/                        # Shared packages
│   ├── domain/                # Domain models
│   ├── errors/                # Standardized error types
│   ├── config/                # Centralized configuration
│   ├── client/                # Keycloak client
│   ├── jaeger/                # Tracing setup
│   ├── middlewares/           # Prometheus metrics middleware
│   └── notification/           # Kafka notification client
├── keycloak/config/           # Keycloak realm configuration
├── scripts/
│   ├── init_sql/              # Database initialization scripts
│   └── generate_from_swagger.sh
├── tests/                      # Integration tests
├── docker-compose.yaml        # Service orchestration
├── Makefile                   # Build automation
└── prometheus.yml             # Prometheus configuration
```

### Architecture Layers

Each service follows a consistent layered architecture:

1. **Handler Layer** (`internal/handlers/` or `internal/restapi/handlers/`)
   - HTTP request/response handling
   - Domain ↔ API model mapping
   - Error conversion
   - Tracing/logging

2. **Service Layer** (`internal/service/` or `internal/database_service/`)
   - Business logic
   - Validation
   - Authorization checks
   - Orchestration

3. **Repository Layer** (`internal/repository/`)
   - Interface-based data access
   - PostgreSQL implementation
   - Query building
   - Error handling

4. **Domain Models** (`pkg/domain/`)
   - Core business entities
   - Validation logic
   - Type-safe enums
   - Separate from API models

## Troubleshooting

### Services Not Starting

If services fail to start:

1. Check logs:
   ```bash
   docker-compose logs <service-name>
   ```

2. Verify `.env` file exists and has all required variables

3. Check if ports are already in use:
   ```bash
   lsof -i :8800  # Auth service
   lsof -i :8880  # Booking service
   lsof -i :8888  # Parking service
   lsof -i :8890  # Payment service
   ```

### Setup Service Fails

If setup fails:

1. Check setup logs:
   ```bash
   docker-compose logs setup
   ```

2. Verify `KEYCLOAK_CLIENT_SECRET` is set in `.env` (not the default value)

3. Restart setup service:
   ```bash
   docker-compose up -d setup
   ```

4. Check if Keycloak is ready:
   ```bash
   docker-compose logs keycloak | tail -20
   ```

### Database Tables Missing

If database tables are missing:

1. Check database logs:
   ```bash
   docker-compose logs db
   ```

2. Manually run setup again:
   ```bash
   docker-compose up -d setup
   docker-compose logs -f setup
   ```

### Keycloak Setup Fails

If Keycloak setup fails:

1. Wait longer - Keycloak can take 30-60 seconds to fully start
2. Check Keycloak logs:
   ```bash
   docker-compose logs keycloak
   ```
3. Verify `KEYCLOAK_CLIENT_SECRET` is set in `.env`
4. Restart setup:
   ```bash
   docker-compose restart setup
   docker-compose logs -f setup
   ```

### Keycloak Configuration Issues

If you see errors like "Invalid client or Invalid client credentials":

1. Check Keycloak is running: `docker-compose ps keycloak`
2. Verify environment variables in `.env`:
   - `KEYCLOAK_CLIENT`
   - `KEYCLOAK_CLIENT_SECRET`
   - `KEYCLOAK_REALM`
   - `KEYCLOAK_ADMIN`
   - `KEYCLOAK_ADMIN_PASSWORD`
3. Ensure Keycloak realm is imported and configured

### Service Not Available

If tests fail with connection errors:

1. Check all services are up: `docker-compose ps`
2. Check service logs: `docker-compose logs <service-name>`
3. Verify ports are not blocked

### Database Issues

If registration/login fails:

1. Check database is healthy: `docker-compose ps db`
2. Verify databases are initialized
3. Check database logs: `docker-compose logs db`

### Payment Service Issues

If payment processing fails:

1. Check payment service logs: `docker-compose logs payment`
2. Verify `INTERNAL_SERVICE_TOKEN` is set in `.env`
3. Check payment database: `docker-compose exec db psql -U postgres -d payment_db -c "SELECT * FROM balances LIMIT 5;"`
4. Verify gRPC connection between booking and payment services

### Common Commands

```bash
# Start all services (setup runs automatically)
docker-compose up -d

# Stop all services
docker-compose down

# Restart all services
docker-compose restart

# View logs
docker-compose logs -f
# or for specific service
docker-compose logs -f <service-name>

# Check service status
docker-compose ps

# Run tests
python3 tests/integration_test.py

# Clean everything (removes volumes)
docker-compose down -v
```

## Security Considerations

- **Authentication**: OAuth2/OIDC via Keycloak
- **Authorization**: Role-based access control (RBAC)
- **Token Validation**: Every protected endpoint validates JWT tokens
- **Database Security**: Connection pooling with pgx driver, parameterized queries
- **Secret Management**: Environment variables (use secret managers in production)
- **Network Isolation**: Services communicate via Docker network
- **Input Validation**: All user inputs are validated
- **SQL Injection Protection**: Parameterized queries throughout
- **Balance Overflow Protection**: Safe arithmetic operations for financial calculations
- **Inter-service Authentication**: gRPC calls authenticated with `INTERNAL_SERVICE_TOKEN`
- **Error Sanitization**: Generic error messages to prevent information disclosure

Security Best Practices:
1. Never commit `.env` file - It's already in `.gitignore`
2. Use strong secrets - Generate random strings for production
3. Rotate secrets regularly in production
4. Use different secrets for development and production environments
5. Enable HTTPS in production (see [HTTPS Setup](#https-setup))
6. Regularly update dependencies
7. Monitor security advisories

## Production Deployment

Before deploying to production:

1. **Configuration Management**: Use proper secret managers (Vault, AWS Secrets Manager)
2. **Database**:
   - Set up replication and backups
   - Use connection pooling
   - Implement migration strategy
   - Regular backup schedule
3. **Monitoring**:
   - Set up alerting rules in Prometheus
   - Configure log aggregation (ELK, Loki)
   - Set up uptime monitoring
   - Configure Grafana dashboards
4. **Scaling**:
   - Configure horizontal pod autoscaling
   - Use load balancers
   - Implement rate limiting
   - Consider service mesh (Istio, Linkerd)
5. **Security**:
   - Enable TLS/SSL for all services (see [HTTPS Setup](#https-setup))
   - Implement API gateway
   - Regular security audits
   - Use WAF (Web Application Firewall)
6. **Performance**:
   - Add caching layer (Redis)
   - Optimize database queries
   - Implement circuit breakers
   - Use CDN for frontend assets
7. **High Availability**:
   - Deploy multiple instances of each service
   - Use database replication
   - Implement health checks
   - Set up automatic failover

## Docker Services

| Service | Container Name | Internal Port | External Port | Purpose |
|---------|---------------|---------------|---------------|---------|
| PostgreSQL | db | 5432 | 5432 | Multi-database persistence |
| Parking | parking-svc | 8888, 50051 | 8888, 50051 | Parking place management |
| Booking | booking-svc | 8880 | 8880 | Booking management |
| Payment | payment-svc | 8890, 50052 | 8890, 50052 | Payment processing |
| Auth | auth-svc | 8800 | 8800 | Authentication |
| Notification | notification-svc | - | - | Notification handling |
| Telegram Bot | telegram | - | - | Bot interface |
| Frontend | parking-frontend | 80 | 3000 | React web interface |
| Keycloak | keycloak | 8080 | 8080 | Identity management |
| Kafka | kafka | 9092 | 9092 | Message broker |
| Zookeeper | zookeeper | 2181 | 2181 | Kafka coordination |
| Jaeger | jaeger | 16686, 14268 | 16686, 14268 | Distributed tracing |
| Prometheus | prometheus | 9090 | 9090 | Metrics collection |
| Grafana | grafana | 3000 | 3000 | Metrics visualization |

## Access Points

- **Services**: Direct access on ports 8800 (auth), 8888 (parking), 8880 (booking), 8890 (payment)
- **For HTTPS**: Configure nginx reverse proxy on host (see [HTTPS Setup](#https-setup))
  - `/auth` - Auth API
  - `/parking` - Parking API
  - `/booking` - Booking API
  - `/payment` - Payment API
- **Frontend** (standalone): http://localhost:3000
- **Keycloak Admin**: http://localhost:8080
- **Jaeger UI**: http://localhost:16686
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
