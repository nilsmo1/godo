BINARY=bin
build:
	go build -o $(BINARY) main.go
run:
	@./$(BINARY)
br: build run 
clean:
	go clean
	rm -f $(BINARY)
