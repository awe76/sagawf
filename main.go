package main

import (
	"github.com/awe76/sagawf/handler"
	pb "github.com/awe76/sagawf/proto"

	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"
)

var (
	service = "sagawf"
	version = "latest"
)

func main() {
	// Create service
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
	)
	srv.Init()

	handler, err := handler.NewSagawf(srv.Client())

	if err != nil {
		log.Fatal(err)
		return
	}

	// Register handler
	pb.RegisterSagawfHandler(srv.Server(), handler)

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
