# Notify Service

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-4285F4?style=for-the-badge&logo=google&logoColor=white)
![NATS](https://img.shields.io/badge/NATS-27AAE1?style=for-the-badge&logo=nats&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)

A high-performance, asynchronous Notification Service built with Go, gRPC, and NATS Jetstream. This service is designed to handle email notifications reliably by leveraging message queuing and a persistent history log.

## 🚀 Features

- **gRPC API**: Fast and typed interface for sending notifications.
- **Asynchronous Processing**: Uses NATS Jetstream to decouple request handling from actual email delivery.
- **Persistent History**: Every notification is logged in MySQL for audit and status tracking.
- **Custom Templates**: Premium HTML email templates for OTP, Account Activation, and Password Reset.
- **Dockerized**: Ready to deploy with Docker and Docker Compose.
- **CI/CD Ready**: Pre-configured GitHub Actions workflow.

## 🛠 Tech Stack

- **Languange**: Go 1.25+
- **Protocol**: gRPC (via Protocol Buffers)
- **Messaging**: NATS Jetstream
- **Database**: MySQL 8.0
- **Configuration**: Viper
- **Containerization**: Docker & Docker Compose
- **Logging**: Standard Library (ready for zap/logrus)

## 📁 Project Structure

```text
.
├── cmd/notify/             # Application entry point
├── internal/
│   ├── entity/             # Domain models (Notification details, Constants)
│   ├── dto/                # Data Transfer Objects
│   ├── usecase/            # Business logic & Background worker
│   ├── repository/         # Database persistence (MySQL)
│   └── handler/grpc/       # gRPC server implementation
├── pkg/
│   ├── mysql/              # Database connection logic
│   ├── nats/               # NATS Jetstream initialization
│   └── email/              # SMTP sending helper
├── proto/notify/           # Protocol Buffer definitions
├── gen/go/notify/          # Generated Go gRPC code
└── templates/              # HTML Email Templates
```

## 🚥 Getting Started

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go 1.25+](https://go.dev/dl/) (optional for local builds)
- [grpcurl](https://github.com/fullstorydev/grpcurl) (for testing)

### Running with Docker Compose

1. Clone the repository
2. Run the stack:
   ```bash
   docker-compose up --build
   ```
3. Access the Mailhog dashboard at [http://localhost:8025](http://localhost:8025) to see sent emails.

## 📡 API Usage

### Send Notification

**Service**: `notify.v2.NotifyService`
**Method**: `Send`

Example request using `grpcurl`:

```bash
grpcurl -plaintext -d '{
  "email": "user@example.com",
  "phone": "08123456789",
  "type": "OTP",
  "data": "123456",
  "metadata": "{\"source\": \"auth-service\"}"
}' localhost:50057 notify.v2.NotifyService/Send
```

### Notification Types

- `OTP`: Sends a security verification code.
- `ACTIVATION`: Sends an account activation token.
- `RESET_PASSWORD`: Sends a password reset link/token.

## ⚙️ Configuration

Environment variables can be set in `docker-compose.yml` or a `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `MYSQL_HOST` | MySQL Server Address | `localhost` |
| `NATS_URL` | NATS Connection String | `nats://localhost:4222` |
| `SMTP_HOST` | SMTP Server Host | `mailhog` |
| `SMTP_PORT` | SMTP Server Port | `1025` |
| `GRPC_PORT` | gRPC Server Port | `50057` |

## 👨‍💻 Author

- **ImamTry257**

---
© 2026 Notify Service. All rights reserved.
