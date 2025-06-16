# Simple URL Compressor

![Build](https://img.shields.io/badge/build-passing-brightgreen.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)

This is a backend project for order and payment microservices, featuring asynchronous communication, reliable event delivery, and basic account management.

## How to Run

### Prerequisites

Make sure you have the following installed on your system:
- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Running the Application

To start the application, run the following command in your terminal:

```bash
POSTGRES_PASSWORD=<your_password> docker-compose up -d
```

Replace `<your_password>` with your desired PostgreSQL password.

---

### API Gateway

We utilize [Kong](https://github.com/Kong/kong) as the API gateway for routing requests. Kong acts as the entry point for all incoming requests, efficiently distributing them to the appropriate services. 

---

### URL Endpoints

# Available Endpoints via Kong

## Order service

- **POST /order** — create a new order
  ```sh
  curl -X POST http://localhost:8000/order \
    -H "Content-Type: application/json" \
    -d '{
      "user_id": "13112111-1111-1111-1111-111111111111",
      "amount": 1000,
      "description": "Test order"
    }'
  ```

- **GET /order/{order_id}/status** — get order status
  ```sh
  curl http://localhost:8000/order/{order_id}/status
  ```

- **GET /order?user_id=...** — get all orders for a user
  ```sh
  curl "http://localhost:8000/order?user_id=13112111-1111-1111-1111-111111111111"
  ```

## Payment service

- **POST /payment/accounts** — create a new account
  ```sh
  curl -X POST http://localhost:8000/payment/accounts \
    -H "Content-Type: application/json" \
    -d '{
      "user_id": "13112111-1111-1111-1111-111111111111"
    }'
  ```

- **GET /payment/accounts/{user_id}/balance** — get user account balance
  ```sh
  curl http://localhost:8000/payment/accounts/13112111-1111-1111-1111-111111111111/balance
  ```

- **POST /payment/accounts/{user_id}/deposit** — deposit to user account
  ```sh
  curl -X POST http://localhost:8000/payment/accounts/13112111-1111-1111-1111-111111111111/deposit \
    -H "Content-Type: application/json" \
    -d '{"amount": 1000}'
  ```

# Architecture Notes

- **Outbox pattern**: Used in the order service to reliably publish events to Kafka.
- **Inbox pattern**: Used in the payment service for idempotent event processing from Kafka.
- **Kafka**: Provides asynchronous, decoupled communication between services.

