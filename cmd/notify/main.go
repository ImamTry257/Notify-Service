package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ImamTry257/Notify-Service/config"
	"github.com/ImamTry257/lms-proto-notify/gen/go/notify"
	handler "github.com/ImamTry257/Notify-Service/internal/handler/grpc"
	"github.com/ImamTry257/Notify-Service/internal/repository"
	"github.com/ImamTry257/Notify-Service/internal/usecase"
	"github.com/ImamTry257/Notify-Service/pkg/email"
	"github.com/ImamTry257/Notify-Service/pkg/mysql"
	"github.com/ImamTry257/Notify-Service/pkg/nats"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. MySQL Connection
	db, err := mysql.NewMySQLConnection(cfg.MySQL)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()
	log.Println("Connected to MySQL successfully")

	// 3. NATS Jetstream Connection
	js, nc, err := nats.NewNATSConnection(cfg.NATS)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()
	log.Println("Connected to NATS Jetstream successfully")

	// 4. Initialize Repository, Usecase, and Handler
	repo := repository.NewMySQLNotifyRepository(db)
	uc := usecase.NewNotifyUsecase(repo, js, cfg.NATS.Subject)
	h := handler.NewNotifyHandler(uc)

	// 5. Initialize Email Sender
	emailSender := email.NewSMTPSender(cfg.SMTP)

	// 6. Start NATS Worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker := usecase.NewNATSWorker(js, repo, emailSender, cfg.NATS.Subject, "notify-service-durable")
	if err := worker.Start(ctx); err != nil {
		log.Fatalf("Failed to start NATS worker: %v", err)
	}
	log.Println("NATS Worker started")

	// 6. Start gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer()
	notify.RegisterNotifyServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	go func() {
		log.Printf("gRPC server listening on port %d", cfg.GRPC.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	grpcServer.GracefulStop()
	cancel()
	log.Println("Shutdown complete")
}
