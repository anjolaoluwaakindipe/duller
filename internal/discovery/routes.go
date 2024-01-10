package discovery

import (
	"encoding/json"
	"net/http"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/gorilla/mux"
)

type MuxRouter struct {
	balancer balancer.LoadBalancer
	registry service.Registry
}

func (rt *MuxRouter) SendHeartBeat() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		var message HeartBeatMessage

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
		}

		if err := rt.registry.RegisterService(service.ServiceInfo{ServiceId: message.ServiceId, Path: message.Path, Port: message.Port, IP: message.Port}); err != nil {
			wr.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (rt *MuxRouter) GetServiceMessage() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		var message GetServiceMessage
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
		}

		serviceInfo, err := rt.balancer.GetNextService(message.Path)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
		}

		json.NewEncoder(wr).Encode(serviceInfo)
	}
}

func (rt *MuxRouter) ShowServices() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
	}
}

func (rt *MuxRouter) SetupRoutes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("heartbeat", rt.SendHeartBeat()).Methods("POST")
	router.HandleFunc("getservice/:path", rt.GetServiceMessage()).Methods("GET")

	return router
}

type Router interface {
	SetupRoutes() http.Handler
}

func CreateMuxRouter(balancer balancer.LoadBalancer, registry service.Registry) *MuxRouter {
	return &MuxRouter{balancer: balancer, registry: registry}
}
