package grpc

import (
	"fmt"
	"log"
	"net"

	notifypbv2 "github.com/ImamTry257/lms-proto-notify/gen/go/notify"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	port string
	s    *grpc.Server
}

func New(port string) *Server {
	return &Server{
		port: port,
		s:    grpc.NewServer(),
	}
}

func (g *Server) RegisterServices(notifyHandler notifypbv2.NotifyServiceServer) {
	notifypbv2.RegisterNotifyServiceServer(g.s, notifyHandler)
	reflection.Register(g.s)
}

func (g *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", g.port))
	if err != nil {
		return err
	}

	log.Printf("gRPC server running on port %s ...", g.port)
	return g.s.Serve(lis)
}

func (g *Server) GracefulStop() {
	log.Println("stopping gRPC server gracefully...")
	g.s.GracefulStop()
	log.Println("gRPC server stopped")
}
