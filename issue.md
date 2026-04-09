# Project Plan: Notify Service (Go)

## Overview
Create a high-level implementation plan for a Go-based Notify Service. This service handles email notifications via gRPC and utilizes NATS Jetstream for asynchronous processing.

## 1. Project Organization
- **Language**: Go
- **Architecture**: Clean Architecture
    - `cmd/`: Entry point (Main)
    - `internal/handler/grpc/`: gRPC server implementation
    - `internal/usecase/`: Business logic
    - `internal/repository/`: Data access (MySQL)
    - `internal/entity/`: Domain models
    - `internal/dto/`: Data transfer objects
    - `config/`: Configuration management
    - `middleware/`: gRPC interceptors (Logging, Recovery)

## 2. External Integration
- **Protocol Buffers**: Import and use `lms-proto-notify` from `git@github.com:ImamTry257/lms-proto-notify.git`.
- **Database**: MySQL for persistence and logging of notification history.
- **Messaging (NATS Jetstream)**:
    - **Stream Name**: `notify`
    - **Subject**: `notify.email`
    - **Pattern**: The service should act as both a producer (accepting requests) and a worker (consuming notifications to send emails).

## 3. Database Schema: `email_histories`
The service should use a table named `email_histories` in MySQL to track all notification attempts:
- `id` (BIGINT, Primary Key, Auto Increment)
- `email` (VARCHAR)
- `phone` (VARCHAR, Optional)
- `type` (VARCHAR) - e.g., 'otp', 'reset_password'
- `data` (TEXT) - Main message content or payload
- `additional_data` (TEXT, Optional)
- `status` (VARCHAR) - e.g., 'PENDING', 'SENT', 'FAILED'
- `metadata` (JSON or TEXT) - For storing additional context
- `sent_at` (TIMESTAMP, NULL)
- `created_at` (TIMESTAMP, Default CURRENT_TIMESTAMP)
- `updated_at` (TIMESTAMP, Default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)

## 4. Data Structure (Notify Service Core)
The gRPC and DTO layers must handle the following fields:
- `email`, `phone`, `type`, `data`, `additional_data`, `status`, `sent_at`, `created_at`, `metadata`.

## 5. Implementation Steps (High-Level)
1. **Repository Setup**: Initialize the Go repository and folder structure.
2. **Proto Generation**: Integrate the `lms-proto-notify` repository and generate Go code using `protoc`.
3. **Database Layer**: Setup MySQL connection and repository methods for saving notification logs.
4. **NATS Integration**: Implement NATS Jetstream connection and stream/consumer management.
5. **Usecase & Logic**: 
    - Implement the logic to publish to NATS when a gRPC request is received.
    - Implement a background worker to consume from NATS and perform the actual notification dispatch (email).
6. **gRPC Interface**: Implement the gRPC handler to provide an API for other services.
7. **Infrastructure**:
    - **Docker**: Create `Dockerfile` and `docker-compose.yml`.
    - **CI/CD**: Setup GitHub Actions for automated build and test pipelines.

## 5. Deployment & Configuration
- Use environment variables for all sensitive and environment-specific configs (DB credentials, NATS URL, SMTP settings).
- Ensure the service is containerized for easy deployment.
