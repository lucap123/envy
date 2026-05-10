build:
	go build -o envy.exe main.go

test:
	go test ./...

clean:
	rm -f envy.exe
	rm -f .env .env.example
