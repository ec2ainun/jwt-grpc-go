package main

import (
	"log"
	"net"
	"os"

	config "github.com/ec2ainun/jwt-grpc-go/config"
	authpb "github.com/ec2ainun/jwt-grpc-go/proto/auth"
	serve "github.com/ec2ainun/jwt-grpc-go/server/auth/serve"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var conf = config.Value()
var err error

func main() {
	listenAddr := conf.AuthServerAddr + ":" + conf.AuthPort
	jwtKey := conf.JWTKey
	tlsCert := conf.TLSCert
	tlsKey := conf.TLSKey
	username := conf.Username
	password := conf.Password

	logger := log.New(os.Stderr, "At ", log.LstdFlags)
	logger.Printf("Starting AuthService...")

	if username == "" || password == "" {
		logger.Fatalln("Please provide Username and Password")
	}

	creds, err := credentials.NewServerTLSFromFile(tlsCert, tlsKey)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	gs := grpc.NewServer(grpc.Creds(creds))

	as, err := serve.NewAuthServer(jwtKey, username, password)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	authpb.RegisterAuthServiceServer(gs, as)

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Could not Listen : %v", err)
	}
	if err := gs.Serve(lis); err != nil {
		logger.Fatalf("Failed to Serve: %v", err)
	}
	logger.Printf("Listening At %s", listenAddr)
}
