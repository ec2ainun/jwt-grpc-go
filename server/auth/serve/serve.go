package serve

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	authpb "github.com/ec2ainun/jwt-grpc-go/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

//Server ...
type Server struct {
	jwtPrivateKey *rsa.PrivateKey
	username      string
	password      string
}

//Login ...
func (s *Server) Login(ctx context.Context, req *authpb.ReqLogin) (*authpb.RespLogin, error) {
	logger := log.New(os.Stderr, "At ", log.LstdFlags)
	logger.Printf("Received Login RPC...")

	if req.Username != s.username || req.Password != s.password {
		return nil, grpc.Errorf(codes.PermissionDenied, "Invalid Username Password")
	}

	now := time.Now()
	exp := now.Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud": "grpc.hello.default",
		"iss": "grpc.auth.default",
		"exp": exp.Unix(),
		"iat": now.Unix(),
		"usr": req.Username,
	})

	tokenStr, err := token.SignedString(s.jwtPrivateKey)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	return &authpb.RespLogin{Token: tokenStr}, nil
}

//NewAuthServer ...
func NewAuthServer(jwtKey, username, password string) (*Server, error) {
	rawJwtKey, err := ioutil.ReadFile(jwtKey)
	if err != nil {
		return nil, fmt.Errorf("Could not Load PrivateKey JWT: %v", err)
	}

	parsedJwtKey, err := jwt.ParseRSAPrivateKeyFromPEM(rawJwtKey)
	if err != nil {
		return nil, fmt.Errorf("Could not Parse PrivateKey JWT: %v", err)
	}

	return &Server{
		jwtPrivateKey: parsedJwtKey,
		username:      username,
		password:      password,
	}, nil
}
