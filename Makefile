client_file := "cmd/client/main.go"

discovery_file := "./cmd/discovery/main.go"

dev:client

client:discovery
	go run $(client_file)

discovery:
	go run $(discovery_file)
