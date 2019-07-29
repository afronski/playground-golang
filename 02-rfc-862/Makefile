GOCMD=go

SERVER_BINARY_NAME=echo-server
CLIENT_BINARY_NAME=echo-client

all: build

build:
	$(GOCMD) build -v -o $(SERVER_BINARY_NAME) echo-server.go
	$(GOCMD) build -v -o $(CLIENT_BINARY_NAME) echo-client.go

clean:
	$(GOCMD) clean
	rm -f $(SERVER_BINARY_NAME)
	rm -f $(CLIENT_BINARY_NAME)

format:
	$(GOCMD) fmt echo-server.go
	$(GOCMD) fmt echo-client.go

vet:
	$(GOCMD) vet echo-server.go
	$(GOCMD) vet echo-client.go

run-server:
	$(GOCMD) run -v ./echo-server.go

run-client:
	$(GOCMD) run -v ./echo-client.go
