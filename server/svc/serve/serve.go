package serve

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	hellopb "github.com/ec2ainun/jwt-grpc-go/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

//Server ...
type Server struct {
	jwtPubKey *rsa.PublicKey
}

type customClaims struct {
	User string `json:"usr"`
	jwt.StandardClaims
}

//Greet ...
func (s *Server) Greet(ctx context.Context, req *hellopb.ReqGreet) (*hellopb.RespGreet, error) {
	logger := log.New(os.Stderr, "At ", log.LstdFlags)
	logger.Printf("Received Greet RPC...")

	claims, err := getClaim(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get Claims JWT : %v", err)
	}
	logger.Printf("In Greet RPC, your Username is: %s", claims.User)

	return &hellopb.RespGreet{
		Message: "Hello " + req.Name + ", your Username is: " + claims.User + ", Welcome!",
	}, nil
}

//PrimeDecompose ...
func (s *Server) PrimeDecompose(req *hellopb.ReqPrimeDecompose, stream hellopb.HelloService_PrimeDecomposeServer) error {
	logger := log.New(os.Stderr, "At ", log.LstdFlags)
	logger.Printf("Received PrimeDecompose RPC...")

	number := req.GetNumber()
	divisor := int64(2)

	claims, err := getClaim(stream.Context(), s)
	if err != nil {
		return fmt.Errorf("Failed to Get Claims JWT : %v", err)
	}
	logger.Printf("In PrimeDecompose RPC, your Username is: %s", claims.User)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&hellopb.RespPrimeDecompose{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
		}
	}
	return nil
}

// ComputeAverage ...
func (s *Server) ComputeAverage(stream hellopb.HelloService_ComputeAverageServer) error {
	logger := log.New(os.Stderr, "At ", log.LstdFlags)
	logger.Printf("Received ComputeAverage RPC...")

	sum := int32(0)
	count := 0

	claims, err := getClaim(stream.Context(), s)
	if err != nil {
		return fmt.Errorf("Failed to Get Claims JWT : %v", err)
	}
	logger.Printf("In ComputeAverage RPC, your Username is: %s", claims.User)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&hellopb.RespComputeAverage{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while Reading Client Stream: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}

//FindMax ...
func (s *Server) FindMax(stream hellopb.HelloService_FindMaxServer) error {
	logger := log.New(os.Stderr, "At ", log.LstdFlags)
	logger.Printf("Received FindMax RPC...")

	claims, err := getClaim(stream.Context(), s)
	if err != nil {
		return fmt.Errorf("Failed to Get Claims JWT : %v", err)
	}
	logger.Printf("In FindMax RPC, your Username is: %s", claims.User)

	maximum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		number := req.GetNumber()
		if number > maximum {
			maximum = number
			sendErr := stream.Send(&hellopb.RespFindMax{
				MaxNumber: maximum,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		}
	}

}

func getClaim(ctx context.Context, s *Server) (*customClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "Valid JWT Token Required")
	}

	jwtToken, ok := md["authorization"]
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "Valid JWT Token Required")
	}

	_, claims, err := validateJwtToken(jwtToken[0], s.jwtPubKey)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "Valid JWT Token Required: %v", err)
	}

	return claims, nil
}

func validateJwtToken(token string, key *rsa.PublicKey) (*jwt.Token, *customClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Valid JWT Token Required")
		}
		return key, nil
	})

	if claims, ok := jwtToken.Claims.(*customClaims); ok && jwtToken.Valid {
		return jwtToken, claims, nil
	}

	return nil, nil, err
}

// NewHelloServer ...
func NewHelloServer(jwtKey string) (*Server, error) {
	rawJwtKey, err := ioutil.ReadFile(jwtKey)
	if err != nil {
		return nil, fmt.Errorf("Could not Load PubKey JWT: %v", err)
	}

	parsedJwtKey, err := jwt.ParseRSAPublicKeyFromPEM(rawJwtKey)
	if err != nil {
		return nil, fmt.Errorf("Could not Parse PubKey JWT: %v", err)
	}

	return &Server{
		jwtPubKey: parsedJwtKey,
	}, nil
}
