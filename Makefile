BINARY=godo-bin
build:
	go build -o $(BINARY) main.go

run:
	@./$(BINARY)

br: build run 

install: build
	chmod +x $(BINARY)
	sudo cp $(BINARY) /usr/bin/
	test -f ~/.config/godo-lists.json || cp template.json ~/.config/godo-lists.json
	go clean
	rm -f $(BINARY)

clean:
	go clean
	rm -f $(BINARY)
