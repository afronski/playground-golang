GOCMD=go

BINARY_NAME=bf-repl

all: build

build:
	$(GOCMD) build -o $(BINARY_NAME) -v

clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)

format:
	$(GOCMD) fmt

vet:
	$(GOCMD) vet

run:
	$(GOCMD) run -v ./bf-repl.go
