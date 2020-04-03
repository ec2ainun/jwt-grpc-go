package main

import (
	"log"
	"net"
	"os"

	config "github.com/ec2ainun/jwt-grpc-go/config"
	hellopb "github.com/ec2ainun/jwt-grpc-go/proto/hello"
	serve "github.com/ec2ainun/jwt-grpc-go/server/svc/serve"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var conf = config.Value()
var err error

func main() {
	listenAddr := conf.HelloServerAddr + ":" + conf.HelloPort
	jwtPubKey := conf.JWTPubKey
	tlsCert := conf.TLSCert
	tlsKey := conf.TLSKey

	logger := log.New(os.Stderr, "At ", log.LstdFlags)
	logger.Printf("Starting HelloService...")

	creds, err := credentials.NewServerTLSFromFile(tlsCert, tlsKey)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	gs := grpc.NewServer(grpc.Creds(creds))

	hs, err := serve.NewHelloServer(jwtPubKey)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	hellopb.RegisterHelloServiceServer(gs, hs)

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Failed to Listen: %v", err)
	}
	if err := gs.Serve(lis); err != nil {
		logger.Fatalf("Failed to Serve: %v", err)
	}
	logger.Printf("Listening At %s", listenAddr)
}
