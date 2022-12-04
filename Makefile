BINARY=godo-bin
build:
	go build -o $(BINARY) main.go

run:
	@./$(BINARY)

br: build run 

install: build
	chmod +x $(BINARY)
	sudo cp $(BINARY) /usr/bin/
	cp template.json ~/.config/godo-lists.json
	go clean

clean:
	go clean
	rm -f $(BINARY)
