client_file := "./cmd/client/main.go"

client_air_file := "./cmd/client/.air.toml"

discovery_file := "./cmd/duller/main.go"

discovery_air_file := "./cmd/duller/.air.toml"

DFLAGS = ""
GFLAGS = ""

dev-client:
	cd ./cmd/client/ && air

dev-disc:
	air -c "./cmd/duller/.air.toml" disc $(DFLAGS)

dev-gate:
	air -c "./cmd/duller/.air.toml" gate $(GFLAGS)
