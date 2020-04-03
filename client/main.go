package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ec2ainun/jwt-grpc-go/config"
	"github.com/ec2ainun/jwt-grpc-go/jwt"
	authpb "github.com/ec2ainun/jwt-grpc-go/proto/auth"
	hellopb "github.com/ec2ainun/jwt-grpc-go/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var conf = config.Value()
var err error
var logger = log.New(os.Stderr, "At ", log.LstdFlags)

func main() {
	dialAddrAuth := conf.AuthServerAddr + ":" + conf.AuthPort
	dialAddrHello := conf.HelloServerAddr + ":" + conf.HelloPort
	caCert := conf.CACrt
	jwtToken := conf.JWTToken
	username := conf.Username
	password := conf.Password

	logger.Printf("Starting Client...")

	tCa, err := credentials.NewClientTLSFromFile(caCert, "")
	if err != nil {
		logger.Fatalf("%v", err)
	}

	connAuth, err := grpc.Dial(
		dialAddrAuth,
		grpc.WithTransportCredentials(tCa),
	)
	if err != nil {
		logger.Fatalf("Failed to Connect : %v", err)
	}
	defer connAuth.Close()

	clientAuth := authpb.NewAuthServiceClient(connAuth)

	req := &authpb.ReqLogin{Username: username, Password: password}
	res, err := clientAuth.Login(context.Background(), req)
	if err != nil {
		logger.Fatalf("Failed to Login : %v", err)
	}
	err = ioutil.WriteFile(jwtToken, []byte(res.Token), 0600)
	if err != nil {
		logger.Fatalf("Failed to Save Token : %v", err)
	}
	logger.Println("Logged in")

	jwtCreds, err := jwt.NewFromTokenFile(jwtToken)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	connHello, err := grpc.Dial(
		dialAddrHello,
		grpc.WithTransportCredentials(tCa),
		grpc.WithPerRPCCredentials(jwtCreds),
	)
	if err != nil {
		logger.Fatalf("Failed to Connect : %v", err)
	}
	defer connHello.Close()
	clientHello := hellopb.NewHelloServiceClient(connHello)

	activateUnary(clientHello)
	// activateServerStream(clientHello)
	// activateClientStream(clientHello)
	// activateBiDiStream(clientHello)
}

func activateUnary(clientHello hellopb.HelloServiceClient) {
	logger.Printf("Starting Greet RPC...")
	respGreet, err := clientHello.Greet(
		context.Background(),
		&hellopb.ReqGreet{
			Name: "Ainun",
		},
	)
	if err != nil {
		logger.Fatalf("Failed to Greet : %v", err)
	}
	logger.Printf("Got Result : %v", respGreet.Message)
}

func activateServerStream(clientHello hellopb.HelloServiceClient) {
	logger.Printf("Starting PrimeDecompose RPC...")
	respPrime, err := clientHello.PrimeDecompose(
		context.Background(),
		&hellopb.ReqPrimeDecompose{
			Number: 123903928406757657,
		},
	)
	if err != nil {
		logger.Fatalf("Failed to PrimeDecompose : %v", err)
	}
	for {
		res, err := respPrime.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		logger.Printf("Got Result : %v", res.PrimeFactor)
	}
}

func activateClientStream(clientHello hellopb.HelloServiceClient) {
	logger.Printf("Starting ComputeAverage RPC...")
	numberToStream := []int32{3, 5, 9, 1, 5, 6}
	streamAverage, err := clientHello.ComputeAverage(context.Background())
	if err != nil {
		logger.Fatalf("Error occurred while Opening Stream : %v", err)
	}

	for _, number := range numberToStream {
		logger.Printf("Sending number: %v\n", number)
		streamAverage.Send(&hellopb.ReqComputeAverage{
			Number: number,
		})
	}
	resAverage, err := streamAverage.CloseAndRecv()
	if err != nil {
		logger.Fatalf("Error while receiving response: %v", err)
	}
	logger.Printf("Got Average Result: %.3v\n", resAverage.Average)
}

func activateBiDiStream(clientHello hellopb.HelloServiceClient) {
	logger.Printf("Starting FindMax RPC...")
	streamFinding, err := clientHello.FindMax(context.Background())
	if err != nil {
		logger.Fatalf("Error occurred while Opening Stream : %v", err)
	}

	waitc := make(chan struct{})
	// send go routine
	go func() {
		numberToStream := []int32{4, 7, 2, 9, 5, 10, 4}
		for _, number := range numberToStream {
			logger.Printf("Sending number: %v\n", number)
			streamFinding.Send(&hellopb.ReqFindMax{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		streamFinding.CloseSend()
	}()
	// receive go routine
	go func() {
		for {
			res, err := streamFinding.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Printf("Error occurred while Reading Server Stream: %v\n", err)
				break
			}
			logger.Printf("Got New Max Number : %v", res.MaxNumber)
		}
		close(waitc)
	}()
	<-waitc

}
