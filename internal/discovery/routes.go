package discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/anjolaoluwaakindipe/duller/internal/tmpl"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type MuxRouter struct {
	balancer      balancer.LoadBalancer
	registry      service.Registry
	upgrader      *websocket.Upgrader
	hashSecretKey string
	hub           *hub
}

func (rt *MuxRouter) SendHeartBeat() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		var message HeartBeatMessage

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(wr).Encode(err)
		}
		newService := &service.ServiceInfo{ServiceId: message.ServiceId, Path: message.Path, Port: message.Port, IP: message.Port}
		if err := rt.balancer.AddService(newService); err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}

		updatedServices := rt.registry.GetServices()
		listComponent := make([]tmpl.Service, 0)
		for _, updatedService := range updatedServices {
			listComponent = append(listComponent, tmpl.Service{Port: updatedService.Port, Path: updatedService.Path, ServiceId: updatedService.ServiceId, IP: updatedService.IP, IsHealthy: updatedService.IsHealthy})
		}
		buffer := new(bytes.Buffer)
		comp := tmpl.ServiceListComponent(listComponent)
		err = comp.Render(context.Background(), buffer)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
		}

		rt.hub.broadcaster <- buffer.Bytes()
	}
}

// GetServiceMessage takes in a request with any http method and utilizes the LoadBalancer
// to proxy the user request to a service instance
func (rt *MuxRouter) GetServiceMessage() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		path := params["path"]

		utils.MakeUrlPathValid(&path)

		serviceInfo, err := rt.balancer.GetNextService(path)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
		}
		serviceLocation := "http://" + serviceInfo.IP + ":" + serviceInfo.Port

		proxy, err := utils.ProxyRequest(serviceLocation)
		if err != nil {
			wr.WriteHeader(http.StatusInternalServerError)
		}

		r.URL.Path = strings.TrimPrefix("/getService", r.URL.Path)

		proxy.ServeHTTP(wr, r)
	}
}

// ShowServices renders a page where all services can be seen
func (rt *MuxRouter) ShowServices() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		services := rt.registry.GetServices()

		serviceVal := make([]tmpl.Service, 0)

		for _, val := range services {
			newSrvComp := tmpl.Service{Port: val.Port, Path: val.Path, IsHealthy: val.IsHealthy, IP: val.IP, ServiceId: val.ServiceId}
			serviceVal = append(serviceVal, newSrvComp)
		}

		page := tmpl.Layout(tmpl.Services(serviceVal))

		page.Render(context.Background(), wr)
	}
}

func (rt *MuxRouter) ServicesSocket() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		conn, err := rt.upgrader.Upgrade(wr, r, nil)
		if err != nil {
			// Handle upgrade error
			fmt.Println("Error upgrading to WebSocket:", err)
			return
		}

		defer conn.Close()

		newClient := newSocketClient(rt.hub, conn, WithWriteWaitTime(10*time.Second))
		rt.hub.register <- &newClient

		newClient.readPipe()
	}
}

func (rt *MuxRouter) GetStaticFiles() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(wr, r)
	}
}

func (rt *MuxRouter) SetupRoutes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", rt.ShowServices()).Methods("GET")
	router.HandleFunc("/heartbeat", rt.SendHeartBeat()).Methods("POST")
	router.HandleFunc("/get-service/{path}", rt.GetServiceMessage())
	router.HandleFunc("/services-socket", rt.ServicesSocket())
	router.PathPrefix("/static/").HandlerFunc(rt.GetStaticFiles())
	return router
}

type Router interface {
	// SetupRoutes returns a handler of already connected routes
	SetupRoutes() http.Handler
}

func WithSecretKey(key string) MuxRouterOptions {
	return func(router *MuxRouter) error {
		if len(strings.Trim(key)) == 0 {
			return nil
		}

		bytes, err := bcrypt.GenerateFromPassword([]byte(key), 10)
		if err != nil {
			return err
		}

		router.hashSecretKey = string(bytes)
		return nil
	}
}

type MuxRouterOptions func(*MuxRouter) error

// New MuxRouter instantiates a MuxRouter with all the necessary handleFuncs utilizing the
// service registry. The MuxRouter implements the Router interface for
func NewMuxRouter(balancer balancer.LoadBalancer, registry service.Registry, opts ...MuxRouterOptions) (Router, error) {
	h := newHub()
	go h.run()
	router := &MuxRouter{
		balancer: balancer,
		registry: registry,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hub: &h,
	}

	for _, opt := range opts {
		err := opt(router)
		if err != nil {
			return nil, err
		}
	}

	return router, nil
}
