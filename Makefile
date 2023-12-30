client_file := "cmd/client/main.go"

discovery_file := "./cmd/discovery/main.go"

dev:client

client:
	go run $(client_file)

disc:
	go run $(discovery_file)
