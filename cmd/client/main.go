package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"

	"andrew528i/wisdom_server/internal"
	"andrew528i/wisdom_server/internal/config"
	"andrew528i/wisdom_server/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// HostDesc is a const with the description for the "-host" flag
	HostDesc = "host to bind socket"

	// PortDesc is a const with the description for the "-port" flag
	PortDesc = "port to bind socket"
)

func main() {
	host := flag.String("host", config.DefaultHost, HostDesc)
	port := flag.Uint("port", config.DefaultPort, PortDesc)
	flag.Parse()

	// target is the gRPC endpoint to connect to
	target := fmt.Sprintf("%s:%d", *host, *port)

	log.WithField("host", *host).
		WithField("port", *port).
		Info("connecting to the server")

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := proto.NewWisdomBookClient(conn)

	challengeProto, err := client.GetChallenge(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	challenge := internal.ChallengeFromMessage(challengeProto)
	challenge.Solve(config.Difficulty)

	log.WithField("difficulty", config.Difficulty).
		WithField("nonce", challengeProto.GetNonce()).
		WithField("signature", hex.EncodeToString(challengeProto.GetSignature())).
		WithField("hash", hex.EncodeToString(challenge.Hash())).
		Info("got challenge from server and solved it")

	// get the quote from server based on the solved challenge
	quote, err := client.GetQuote(context.Background(), challenge.Message())
	if err != nil {
		log.Fatal(err)
	}

	log.WithField("quote", quote.GetQuote()).
		Info("got quote from the server")
}
