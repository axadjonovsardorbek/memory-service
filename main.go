package main

import (
	"log"
	cf "memory/config"
	mp "memory/genproto"
	"memory/storage/postgres"
	"memory/service"
	"net"

	"google.golang.org/grpc"
)

func main() {
	config := cf.Load()

	db, err := postgres.NewPostgresStorage(config)

	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", config.MEMORY_SERVICE_PORT)

	if err != nil {
		log.Fatalf("Failed to listen tcp: %v", err)
	}

	s := grpc.NewServer()

	mp.RegisterCommentsServiceServer(s, service.NewCommentsService(db))
	mp.RegisterMediasServiceServer(s, service.NewMediasService(db))
	mp.RegisterMemoriesServiceServer(s, service.NewMemoriesService(db))
	mp.RegisterSharedMemoriesServiceServer(s, service.NewSharedMemoriesService(db))

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
