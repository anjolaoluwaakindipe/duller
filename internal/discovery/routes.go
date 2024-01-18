package discovery

import (
	"encoding/json"
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
		}

		if err := rt.registry.RegisterService(service.ServiceInfo{ServiceId: message.ServiceId, Path: message.Path, Port: message.Port, IP: message.Port}); err != nil {
			wr.WriteHeader(http.StatusInternalServerError)
		}
	}
}

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

func (rt *MuxRouter) ShowServices() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
	}
}

func (rt *MuxRouter) ServicesSocket() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		conn, _ := rt.upgrader.Upgrade(wr, r, nil)

		newClient := socketClient{conn: conn}
		rt.hub.register <- &newClient
	}
}

func (rt *MuxRouter) SetupRoutes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/services", rt.ShowServices()).Methods("GET")
	router.HandleFunc("/heartbeat", rt.SendHeartBeat()).Methods("POST")
	router.HandleFunc("/getService/{path}", rt.GetServiceMessage())

	return router
}

type Router interface {
	SetupRoutes() http.Handler
}

func CreateMuxRouter(balancer balancer.LoadBalancer, registry service.Registry) Router {
	return &MuxRouter{
		balancer: balancer, registry: registry, upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hub: &hub{socketClients: make(map[*socketClient]bool), broadcaster: make(chan []byte)},
	}
}

type socketClient struct {
	conn *websocket.Conn
	send chan []byte
}

type hub struct {
	socketClients map[*socketClient]bool
	broadcaster   chan []byte
	register      chan *socketClient
	unregister    chan *socketClient
}

func (h *hub) removeClient(client *socketClient) {
	delete(h.socketClients, client)
	close(client.send)
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.socketClients[client] = true
		case client := <-h.unregister:
			h.removeClient(client)
		case message := <-h.broadcaster:
			for client := range h.socketClients {
				select {
				case client.send <- message:
				default:
					h.removeClient(client)
				}
			}
		}
	}
}
