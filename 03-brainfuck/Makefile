GOCMD=go

BINARY_NAME=brainfuck

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
	$(GOCMD) run -v ./brainfuck.go
