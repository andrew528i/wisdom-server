package main

import (
	"flag"
	"fmt"
	"net"

	"andrew528i/wisdom_server/internal"
	"andrew528i/wisdom_server/internal/config"
	"andrew528i/wisdom_server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	log "github.com/sirupsen/logrus"
)

const (
	HostDesc = "host to bind socket"
	PortDesc = "port to bind socket"
)

func main() {
	host := flag.String("host", config.DefaultHost, HostDesc)
	port := flag.Uint("port", config.DefaultPort, PortDesc)
	flag.Parse()

	address := fmt.Sprintf("%s:%d", *host, *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	wisdomSrv, err := internal.NewWisdomBookServer()
	if err != nil {
		log.Fatal(err)
	}

	proto.RegisterWisdomBookServer(srv, wisdomSrv)
	reflection.Register(srv)

	log.WithField("host", *host).
		WithField("port", *port).
		Info("running gRPC server")

	if err = srv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
