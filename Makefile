BINARY=envy.exe

build:
	go build -o $(BINARY) main.go

test:
	go test ./...

clean:
ifeq ($(OS),Windows_NT)
	if exist $(BINARY) del $(BINARY)
else
	rm -f $(BINARY)
endif
