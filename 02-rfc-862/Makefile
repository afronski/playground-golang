GOCMD=go

SERVER_BINARY_NAME=echo-server
CLIENT_BINARY_NAME=echo-client

all: build

build:
	$(GOCMD) build -o $(SERVER_BINARY_NAME) -v
	$(GOCMD) build -o $(CLIENT_BINARY_NAME) -v

clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)

format:
	$(GOCMD) fmt

vet:
	$(GOCMD) vet

run-server:
	$(GOCMD) run -v ./echo-server.go

run-client:
	$(GOCMD) run -v ./echo-client.go
