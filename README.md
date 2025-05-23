🛒 E-Commerce Platform with gRPC, NATS, MongoDB, Redis, and Transactional Email
Overview
This project is a scalable e-commerce platform built with a microservices architecture using Go (Golang). It leverages gRPC for REST-to-service communication, NATS for asynchronous messaging between services, MongoDB for persistent storage with transaction support, Redis for caching, and an SMTP-based mailer service for transactional emails. The platform follows Clean Architecture and Domain-Driven Design (DDD) principles to ensure maintainability and scalability.
Key Features

API Gateway: Converts RESTful HTTP requests to gRPC calls for seamless interaction with microservices.
Microservices: User management, inventory management, order processing, and email notifications.
NATS Messaging: Asynchronous communication between Order Service and Inventory Service for stock updates.
MongoDB Transactions: Ensures atomic operations for critical business logic (e.g., order creation).
Redis Caching: Enhances performance by caching frequently accessed data.
Authentication: JWT-based authentication for protected routes.
Transactional Emails: Sends notifications (e.g., order confirmations) via an SMTP-based mailer service.


🧱 Microservices
The platform consists of the following microservices:
1. API Gateway

Purpose: Acts as the entry point, converting REST API requests to gRPC calls and routing them to appropriate services.
Endpoints:
User: /api/users/register, /api/users/authenticate, /api/users/:id
Inventory: /api/inventory, /api/inventory/:id, etc.
Orders: /api/orders, /api/orders/:id, /api/orders/:id/status, etc.
Mailer: /api/email (for sending transactional emails)


Technologies:
Gin: Lightweight web framework for REST API handling.
gRPC: Communicates with backend services.
Middleware: Logging, telemetry, and JWT authentication.



2. User Service

Purpose: Manages user registration, authentication, and profile data.
Features:
Register users with email and password.
Authenticate users and issue JWT tokens.
Retrieve and update user profiles.


gRPC API: Defined in user.proto.
Technologies:
MongoDB for persistent storage.
Redis for session/token caching.
Bcrypt (assumed) for password hashing.



3. Inventory Service

Purpose: Manages product inventory.
Features:
Create, read, update, and delete (CRUD) products.
List products with pagination and category filtering.
Updates stock levels in response to NATS messages from the Order Service.


gRPC API: Defined in inventory.proto.
NATS Integration: Subscribes to NATS topics (e.g., order.created) to process stock updates.
Technologies: MongoDB for storage, Redis for caching, NATS for messaging.

4. Order Service

Purpose: Handles order creation, retrieval, and status updates.
Features:
Create orders with user ID, total, and product items.
Retrieve order details by ID.
Update order status (PENDING, COMPLETED, CANCELLED).
List user orders.
Publishes events to NATS for inventory updates.


gRPC API: Defined in order.proto.
NATS Integration: Publishes messages (e.g., order.created) to notify the Inventory Service of stock changes.
Technologies:
MongoDB with transactions for atomic order creation.
Redis for caching order data or session tokens.
NATS for event-driven communication.



5. Mailer Service

Purpose: Sends transactional emails (e.g., order confirmations, welcome emails).
Features:
Sends emails via SMTP with configurable credentials.


gRPC API: Defined in mailer.proto.
Technologies: Go net/smtp for email delivery, configurable via environment variables.


📡 NATS Integration
NATS is used for asynchronous, event-driven communication between the Order Service and Inventory Service. This ensures decoupled and reliable stock updates when orders are created or updated.
How It Works

Order Service:
When an order is created (via CreateOrder), it publishes a message to a NATS topic (e.g., order.created) containing order details (e.g., product IDs and quantities).
Example message payload:{
  "order_id": "68016c6489e4500884e83831",
  "items": [
    {"product_id": "68016c6489e4500884e8382f", "quantity": 2},
    {"product_id": "68016c6489e4500884e83830", "quantity": 1}
  ]
}




Inventory Service:
Subscribes to the order.created topic.
Processes the message to decrement stock levels in MongoDB, ensuring atomic updates via transactions.


Reliability: NATS ensures messages are delivered even if the Inventory Service is temporarily unavailable, using its built-in message queueing.

Example Implementation
Order Service (Publish)
func (u *OrderUsecase) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
    session, err := u.mongoClient.StartSession()
    if err != nil {
        return "", err
    }
    defer session.EndSession(ctx)

    id, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
        orderID, err := u.repo.CreateOrder(sessionCtx, order)
        if err != nil {
            return nil, err
        }

        // Publish to NATS
        msg := &model.OrderEvent{
            OrderID: orderID,
            Items:   order.Products,
        }
        msgBytes, _ := json.Marshal(msg)
        err = u.natsConn.Publish("order.created", msgBytes)
        if err != nil {
            return nil, err
        }

        return orderID, nil
    })
    if err != nil {
        return "", err
    }
    return id.(string), nil
}

Inventory Service (Subscribe)
func (u *InventoryUsecase) SubscribeToOrders(natsConn *nats.Conn) {
    natsConn.Subscribe("order.created", func(msg *nats.Msg) {
        var event model.OrderEvent
        if err := json.Unmarshal(msg.Data, &event); err != nil {
            log.Printf("Failed to unmarshal order event: %v", err)
            return
        }

        session, err := u.mongoClient.StartSession()
        if err != nil {
            log.Printf("Failed to start session: %v", err)
            return
        }
        defer session.EndSession(context.Background())

        err = session.WithTransaction(context.Background(), func(sessionCtx mongo.SessionContext) (interface{}, error) {
            for _, item := range event.Items {
                if err := u.repo.DecrementStock(sessionCtx, item.ProductID, item.Quantity); err != nil {
                    return nil, err
                }
            }
            return nil, nil
        })
        if err != nil {
            log.Printf("Failed to update stock: %v", err)
        }
    })
}


🧪 Architecture & Design Patterns

Clean Architecture: Separates concerns into layers (handler, usecase, repository, model).
Domain-Driven Design (DDD): Models business domains (users, orders, inventory) with clear boundaries.
gRPC: Enables type-safe, high-performance communication via the API Gateway.
NATS: Provides asynchronous messaging for event-driven workflows.
REST to gRPC Gateway: Translates REST requests to gRPC calls.
MongoDB Transactions: Ensures atomicity for operations like order creation and stock updates.
Redis Caching: Reduces database load by caching user profiles, orders, or tokens.
Middleware: Implements logging, telemetry, and authentication in the API Gateway.
Dependency Injection: Used in handlers and usecases for testability.


⚙️ Installation & Setup
1. Clone the Repository
git clone https://github.com/your-username/ecommerce-platform.git
cd ecommerce-platform

2. Prerequisites

Go: 1.20 or higher
Docker and Docker Compose: For containerized deployment
MongoDB: 5.0 or higher
Redis: 6.0 or higher
NATS: 2.9 or higher
SMTP Server: E.g., Mailtrap, Gmail, or any SMTP provider
protoc: For generating gRPC code from .proto files

3. Install Dependencies
# Install Go dependencies
cd api-gateway
go mod tidy

# Install NATS client
go get github.com/nats-io/nats.go
install it locally also to your os

# Install protoc-gen-go and protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

4. Environment Variables
Create a .env file in the api-gateway directory (or respective service directories):
# API Gateway
PORT=8080
JWT_SECRET=your_jwt_secret_key
INVENTORY_SERVICE=localhost:50051
ORDER_SERVICE=localhost:50052
USER_SERVICE=localhost:50053
MAILER_SERVICE=localhost:50054

# MongoDB (for services)
MONGO_URI=mongodb://localhost:27017

# Redis (for services)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# NATS (for order and inventory services)
NATS_URL=nats://localhost:4222

# SMTP (for mailer service)
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=587
SMTP_USERNAME=your_smtp_username
SMTP_PASSWORD=your_smtp_password

5. Generate gRPC Code
Generate Go code from .proto files:
cd proto
protoc --go_out=../internal/pb --go-grpc_out=../internal/pb *.proto

6. Run Services
Option A: Using Docker Compose
Create a docker-compose.yml (example below) and run:
docker-compose up --build

Option B: Run Locally
# Start NATS server
docker run -d --name nats -p 4222:4222 nats:latest

# API Gateway
cd api-gateway
go run cmd/main.go

# User Service
cd user-service
go run cmd/main.go

# Inventory Service
cd inventory-service
go run cmd/main.go

# Order Service
cd order-service
go run cmd/main.go

# Mailer Service
cd mailer-service
go run cmd/main.go

Example docker-compose.yml
version: '3.8'
services:
  api-gateway:
    build: ./api-gateway
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - JWT_SECRET=your_jwt_secret_key
      - INVENTORY_SERVICE=inventory-service:50051
      - ORDER_SERVICE=order-service:50052
      - USER_SERVICE=user-service:50053
      - MAILER_SERVICE=mailer-service:50054
    depends_on:
      - inventory-service
      - order-service
      - user-service
      - mailer-service
      - mongodb
      - redis
      - nats

  user-service:
    build: ./user-service
    environment:
      - MONGO_URI=mongodb://mongodb:27017
      - MONGO_DB_NAME=ecommerce
      - REDIS_ADDR=redis:6379
    depends_on:
      - mongodb
      - redis

  inventory-service:
    build: ./inventory-service
    environment:
      - MONGO_URI=mongodb://mongodb:27017
      - MONGO_DB_NAME=ecommerce
      - REDIS_ADDR=redis:6379
      - NATS_URL=nats://nats:4222
    depends_on:
      - mongodb
      - redis
      - nats

  order-service:
    build: ./order-service
    environment:
      - MONGO_URI=mongodb://mongodb:27017
      - MONGO_DB_NAME=ecommerce
      - REDIS_ADDR=redis:6379
      - NATS_URL=nats://nats:4222
    depends_on:
      - mongodb
      - redis
      - nats

  mailer-service:
    build: ./mailer-service
    environment:
      - SMTP_HOST=smtp.mailtrap.io
      - SMTP_PORT=587
      - SMTP_USERNAME=your_smtp_username
      - SMTP_PASSWORD=your_smtp_password

  mongodb:
    image: mongo:5.0
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db

  redis:
    image: redis:6.2
    ports:
      - "6379:6379"

  nats:
    image: nats:latest
    ports:
      - "4222:4222"

volumes:
  mongodb_data:


🧪 Testing
Unit Tests
Run unit tests for all services:
cd <service-directory>
go test ./... -v

Mocks are used for MongoDB, Redis, NATS, and gRPC clients using testify/mock.
Integration Tests
Integration tests use Dockerized MongoDB, Redis, and NATS instances. Run:
cd <service-directory>
go test -tags=integration ./...


🔁 Transaction Handling
The Order Service and Inventory Service use MongoDB ACID transactions to ensure atomicity for operations like:

Creating an order and updating inventory stock.
Logging order actions (e.g., status changes).

Example transaction in the Order Service:
func (u *OrderUsecase) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
    session, err := u.mongoClient.StartSession()
    if err != nil {
        return "", err
    }
    defer session.EndSession(ctx)

    id, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
        orderID, err := u.repo.CreateOrder(sessionCtx, order)
        if err != nil {
            return nil, err
        }
        return orderID, nil
    })
    if err != nil {
        return "", err
    }
    return id.(string), nil
}

The Inventory Service processes NATS messages within transactions to ensure consistent stock updates.

💡 Redis Caching
Redis is used for:

Caching: User profiles, product details, or order summaries to reduce database load.
Session Management: Storing JWT tokens or session data.
Rate Limiting: Can be implemented for API endpoints (planned).

Commands for WSL(ubuntu) to check redis cache:
```bash
redis-cli (toLaunch)
keys * (to check) **TTL IS 5 MINUTES**
```

Example usage in user-service:
func (r *UserRepository) GetUserProfile(ctx context.Context, id string) (*model.User, error) {
    cached, err := r.redis.Get(ctx, "user:"+id).Result()
    if err == nil {
        var user model.User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }
    user, err := r.mongo.FindUserByID(ctx, id)
    if err != nil {
        return nil, err
    }
    userJSON, _ := json.Marshal(user)
    r.redis.Set(ctx, "user:"+id, userJSON, 1*time.Hour)
    return user, nil
}


🔐 Security

JWT Authentication: Protected routes in the API Gateway require a valid JWT token.
Password Hashing: User passwords are hashed (assumed bcrypt or scrypt) before storage.
gRPC Security: Uses insecure credentials for development. In production, use TLS.
NATS Security: Configure NATS with TLS and authentication in production.
Input Validation: Handlers validate inputs (e.g., order status, product IDs) to prevent invalid data.


📁 Project Structure
ecommerce-platform/
├── api-gateway/
│   ├── cmd/
│   │   └── main.go
│   ├── config/
│   │   └── config.go
│   ├── internal/
│   │   ├── handler/        # REST handlers and gRPC clients
│   │   ├── middleware/     # Logging, telemetry, auth
│   │   ├── pb/            # Generated gRPC code
│   ├── proto/
│   │   ├── user.proto
│   │   ├── inventory.proto
│   │   ├── order.proto
│   │   ├── mailer.proto
│   └── go.mod
├── user-service/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handler/        # gRPC handlers
│   │   ├── usecase/       # Business logic
│   │   ├── repository/    # MongoDB/Redis access
│   │   ├── model/         # Domain models
│   │   ├── pb/            # Generated gRPC code
│   ├── proto/
│   │   └── user.proto
│   └── go.mod
├── inventory-service/
│   # Similar structure
├── order-service/
│   # Similar structure
├── mailer-service/
│   # Similar structure
├── docker-compose.yml
└── README.md


📡 Protobuf & gRPC
Protobuf files define service contracts:
mailer.proto
syntax = "proto3";
package mailer;
option go_package = "../internal/pb/mailer";
service MailerService {
  rpc SendEmail (SendEmailRequest) returns (SendEmailResponse);
}
message SendEmailRequest {
  string to_email = 1;
  string subject = 2;
  string html_body = 3;
}
message SendEmailResponse {
  string status = 1;
  string message = 2;
}

Generate gRPC code:
protoc --go_out=. --go-grpc_out=. proto/*.proto


📌 Tools & Libraries



Tool/Library
Purpose



Go
Programming language


gRPC
Inter-service communication (API Gateway)


NATS
Asynchronous messaging (Order/Inventory)


Protobuf
Schema definition for gRPC


Gin
REST API framework (API Gateway)


MongoDB
NoSQL database with transactions


Redis
Caching and session management


go-redis
Redis client library


nats.go
NATS client library


net/smtp
SMTP email delivery


testify
Testing and mocking framework


Docker
Containerization


Docker Compose
Multi-container orchestration



📈 Future Enhancements

OpenTelemetry: Add tracing and metrics for observability.
OAuth2: Replace basic JWT with OAuth2 for enhanced security.
GraphQL Gateway: Support GraphQL alongside REST.
Kubernetes: Deploy services with Kubernetes and Helm.
CI/CD: Implement pipelines with GitHub Actions or GitLab CI.
Rate Limiting: Add Redis-based rate limiting for API endpoints.
Email Templates: Use HTML templates for richer email content.
NATS JetStream: Enable persistent messaging for higher reliability.


🤝 Contributing
Contributions are welcome! Please:

Fork the repository.
Create a feature branch (git checkout -b feature/your-feature).
Commit changes (git commit -m 'Add your feature').
Push to the branch (git push origin feature/your-feature).
Open a pull request.

Ensure code passes tests and follows the project structure.

