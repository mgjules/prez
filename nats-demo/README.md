# NATS Demo

A Go demonstration of event-driven modular monolith architecture using [NATS](https://nats.io/) messaging system.

## Overview

This demo showcases a simple e-commerce system with four modules within a single application communicating via NATS:

- **User Module**: Manages user creation and retrieval
- **Product Module**: Handles product inventory and stock management
- **Order Module**: Processes order creation with validation
- **Payment Module**: Simulates payment processing

The system demonstrates both request-reply and publish-subscribe messaging patterns using an embedded NATS server for internal module communication.

## Architecture

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌─────────────────┐
│ User Module │    │Product Module│    │Order Module │    │Payment Module   │
└─────────────┘    └──────────────┘    └─────────────┘    └─────────────────┘
       │                   │                   │                     │
       └───────────────────┼───────────────────┼─────────────────────┘
                           │                   │
                    ┌──────────────────────────────┐
                    │       NATS Server            │
                    │  (In-Process, No TCP)        │
                    └──────────────────────────────┘
```

### Message Patterns

**Request-Reply:**

- `user.create` / `user.get` - User operations
- `product.create` / `product.get` - Product operations
- `order.create` / `order.get` - Order operations

**Publish-Subscribe Events:**

- `events.user.created` - User creation events
- `events.product.created` - Product creation events
- `events.product.stock.updated` - Stock update events
- `events.order.created` - Order creation events
- `events.payment.success` - Payment completion events

## Requirements

- Go 1.23.0 or later

## Installation

```bash
go mod download
```

## Usage

Run the demo:

```bash
go run .
```

The application will:

1. Start an embedded NATS server
2. Initialize all modules
3. Create a user ("John Doe")
4. Add a product ("Widget" with 100 stock)
5. Place an order (10 widgets)
6. Process payment automatically
7. Update product stock

### Example Output

```bash
2025/06/16 17:35:50 INFO nats server started
2025/06/16 17:35:50 INFO connected to nats server
2025/06/16 17:35:50 INFO user created user={ID:ababa8ca-a693-4475-aac5-467ae50a707d}
2025/06/16 17:35:50 INFO product added product="{ID:fc979c80-4de4-4094-acd6-2fadfeaacb83 Name:Widget Stock:100}"
2025/06/16 17:35:50 INFO order placed order="{ID:34aa54f1-79af-4d50-a1e6-6303897fcd80 UserID:ababa8ca-a693-4475-aac5-467ae50a707d ProductID:fc979c80-4de4-4094-acd6-2fadfeaacb83 Quantity:10}"
2025/06/16 17:35:51 INFO payment successful payment="{ID:b0ad615c-05dd-4ad2-a5d7-9a3eaa0562bf OrderID:34aa54f1-79af-4d50-a1e6-6303897fcd80}"
2025/06/16 17:35:51 INFO product stock updated product="{ID:fc979c80-4de4-4094-acd6-2fadfeaacb83 Name:Widget Stock:90}"
```

## Project Structure

```
.
├── main.go     # Application entry point and demo workflow
├── user.go     # User module implementation
├── product.go  # Product module with inventory management
├── order.go    # Order module with validation
├── payment.go  # Payment module with event handling
├── go.mod      # Go module dependencies
└── go.sum      # Dependency checksums
```

## Key Features

- **Embedded NATS Server**: No external dependencies required
- **In-Process Messaging**: Zero TCP overhead for maximum performance
- **Event-Driven Architecture**: Loose coupling between modules
- **Concurrent Operations**: Thread-safe module implementations
- **Validation**: Request validation across all modules
- **Error Handling**: Comprehensive error logging with structured logging

## Dependencies

- [`github.com/nats-io/nats-server/v2`](https://pkg.go.dev/github.com/nats-io/nats-server/v2) - NATS server
- [`github.com/nats-io/nats.go`](https://pkg.go.dev/github.com/nats-io/nats.go) - Go client for NATS
- [`github.com/google/uuid`](https://pkg.go.dev/github.com/google/uuid) - UUID generation

## Message Flow Example

1. **User Creation**: `user.create` → `events.user.created`
2. **Product Creation**: `product.create` → `events.product.created`
3. **Order Placement**: `order.create` → validates user/product → `events.order.created`
4. **Payment Processing**: listens to `events.order.created` → processes payment → `events.payment.success`
5. **Stock Update**: listens to `events.payment.success` → updates inventory → `events.product.stock.updated`

## Learning Outcomes

This demo illustrates:

- NATS request-reply vs publish-subscribe patterns
- Event-driven modular monolith communication
- Module validation and error handling
- In-memory data management with concurrent access
- Structured logging with Go's `slog` package

## (Bonus) Challenge

In some cases, the NATS connection closes before every message is processed.

Example:

```bash
2025/06/16 17:35:43 ERROR failed to publish product stock updated event err="nats: connection closed" product="{ID:a929d226-2bba-4778-bed5-cfd2dc3ba9f8 Name:Widget Stock:90}"
```

How can we fix this? :thinking:

## License

This is a demonstration project for educational purposes.
