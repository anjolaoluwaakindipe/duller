package discovery

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type MuxRouter struct {
	balancer balancer.LoadBalancer
	registry service.Registry
	upgrader websocket.Upgrader
	hub      *hub
}

func (rt *MuxRouter) SendHeartBeat() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		var message HeartBeatMessage

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(wr).Encode(err)
		}

		if err := rt.balancer.AddService(&service.ServiceInfo{ServiceId: message.ServiceId, Path: message.Path, Port: message.Port, IP: message.Port}); err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}

		services, err := json.Marshal(rt.registry.GetServices())
		if err != nil {
			return
		}

		rt.hub.broadcaster <- services
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
		tmpl, err := template.ParseFiles("./templates/services.html")
		if err != nil {
			wr.WriteHeader(http.StatusInternalServerError)
			return
		}

		services := rt.registry.GetServices()

		tmpl.Execute(wr, services)
	}
}

func (rt *MuxRouter) ServicesSocket() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		conn, _ := rt.upgrader.Upgrade(wr, r, nil)

		newClient := newSocketClient(*rt.hub, conn)
		newClient.hub.register <- &newClient

		go newClient.readPipe()
	}
}

func (rt *MuxRouter) SetupRoutes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", rt.ShowServices()).Methods("GET")
	router.HandleFunc("/heartbeat", rt.SendHeartBeat()).Methods("POST")
	router.HandleFunc("/get-service/{path}", rt.GetServiceMessage())
	router.HandleFunc("/services-socket", rt.ServicesSocket())
	fmt.Println("hello")
	return router
}

type Router interface {
	// SetupRoutes returns a handler of already connected routes
	SetupRoutes() http.Handler
}

// New MuxRouter instantiates a MuxRouter with all the necessary handleFuncs utilizing the
// service registry. The MuxRouter implements the Router interface for
func NewMuxRouter(balancer balancer.LoadBalancer, registry service.Registry) Router {
	h := newHub()
	go h.run()
	return &MuxRouter{
		balancer: balancer, registry: registry, upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hub: &h,
	}
}
