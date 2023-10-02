package main

import (
	"flag"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
)

// cli constants
const (
	REGISTRY_HOST = "localhost"
	REGISTRY_PORT = "9876"
	REGISTRY_TYPE = "tcp"
)

func main() {
	rport := flag.String("rport", REGISTRY_PORT, "Registry Port")
	rhost := flag.String("rhost", REGISTRY_HOST, "Registry Host")
	flag.Parse()

	serviceRegistry := discovery.Registry{}

	discovery.InitRegistryServer(discovery.RegistrySettings{
		REGISTRY_HOST: *rhost,
		REGISTRY_PORT: *rport,
		REGISTRY_TYPE: REGISTRY_TYPE,
	}, &serviceRegistry)
}
