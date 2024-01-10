package discovery

import (
	"encoding/json"
	"net/http"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
)

type Routes struct {
	registry service.Registry
	balancer balancer.LoadBalancer
}

func (rt *Routes) RegisterMessage() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		var message RegisterServiceMessage

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
		}

		if err := rt.registry.RegisterService(service.ServiceInfo{ServiceId: message.ServiceName, Path: message.Path}); err != nil {
			wr.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (rt *Routes) GetServicesMessageMessage() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		var message GetAddressMessage
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
		}

		serviceInfo, err := rt.balancer.GetNextService(message.Path)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
		}

		dataMap := make(map[string]interface{})
		dataMap["address"] = "http://" + serviceInfo.IP + ":" + serviceInfo.Port
		response, _ := json.Marshal(dataMap)
		wr.Write(response)
	}
}

func (rt *Routes) ShowServicesMessage() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
	}
}

func messageHandler(message string) {
}

func InitRoutes() *Routes {
	return &Routes{}
}
