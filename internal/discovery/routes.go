package discovery

import (
	"net/http"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
)

type Routes struct {
	registry service.Registry
	balancer balancer.LoadBalancer
}

func (r *Routes) RegisterMessage() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
	}
}

func (r *Routes) GetServiceMessageMessage() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
	}
}

func (r *Routes) ShowServicesMessage() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
	}
}

func InitRoutes() *Routes {
	return &Routes{}
}
