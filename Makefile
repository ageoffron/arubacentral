build:
	go build -o bin/arubacentral_authclient main.go

run:
	go run main.go

compile:
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
