.PHONY: all 
all: gen-cert

gen-cert :
	cd ssl/ && bash generate-certs && cd ../
init-go :
	go mod init
	go mod tidy
	go mod vendor
start-auth-svc :
	go run server/auth/main.go
start-server-svc :
	go run server/svc/main.go
start-client-svc :
	go run client/main.go
