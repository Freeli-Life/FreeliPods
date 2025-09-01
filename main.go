package main

import (
	"crypto"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"FreeliPods/config"
	fc "FreeliPods/crypto"
	"FreeliPods/database"
	pb "FreeliPods/podServer"
	"FreeliPods/server"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	privateKey, err := fc.LoadPrivateKey(cfg.TLSKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	db, err := database.NewStore("users.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	go startGrpcServer(cfg, db, privateKey)

	go startGrpcWebServer(cfg, db, privateKey)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")
}

func startGrpcServer(cfg *config.Config, db *database.Store, key crypto.PrivateKey) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", cfg.Port, err)
	}

	creds, err := credentials.NewServerTLSFromFile(cfg.TLSCert, cfg.TLSKey)
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	podServer := &server.PodServiceServer{DB: db, PrivateKey: key, Domain: cfg.Domain}
	pb.RegisterPodServiceServer(s, podServer)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startGrpcWebServer(cfg *config.Config, db *database.Store, key crypto.PrivateKey) {
	grpcServer := grpc.NewServer()
	podServer := &server.PodServiceServer{DB: db, PrivateKey: key, Domain: cfg.Domain}
	pb.RegisterPodServiceServer(grpcServer, podServer)

	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.WebPort),
		Handler: wrappedGrpc,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("gRPC-Web server listening at https://localhost:%d", cfg.WebPort)
	if err := httpServer.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey); err != nil {
		log.Fatalf("Failed to serve gRPC-Web: %v", err)
	}
}