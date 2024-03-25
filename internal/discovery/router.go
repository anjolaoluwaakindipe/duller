package discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/registry"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/anjolaoluwaakindipe/duller/internal/tmpl"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type MuxRouter struct {
	balancer      balancer.LoadBalancer
	registry      registry.Registry
	upgrader      *websocket.Upgrader
	hashSecretKey string
	ctx           context.Context
	hub           Hub
}

func (rt *MuxRouter) isHeartbeatAuthorized(serviceId string, path string, key string) error {
	err := bcrypt.CompareHashAndPassword([]byte(rt.hashSecretKey), []byte(key))

	if len(rt.hashSecretKey) == 0 {
		return nil
	}

	if err != nil {
		rt.registry.DeregisterService(path, serviceId)
		return fmt.Errorf("Unauthorized Request")
	}
	return nil
}

// getAuthToken returns the token from auth header and checks if the format is correct.
// if the format is wrong or the authorization header is empty an empty string is returned else
// the token is returned
func (rt *MuxRouter) getAuthToken(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(authHeader) == 0 {
		return ""
	}
	authSplit := strings.Split(authHeader, "Bearer ")

	if len(authSplit) != 2 {
		return ""
	}

	return authSplit[1]
}

func (rt *MuxRouter) SendHeartBeat() func(wr http.ResponseWriter, r *http.Request) {
	return func(wr http.ResponseWriter, r *http.Request) {
		var message HeartBeatMessage

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(wr).Encode(err)
		}

		secret := rt.getAuthToken(r)

		if err := rt.isHeartbeatAuthorized(message.ServiceId, message.Path, secret); err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
			return
		}

		newService := &service.ServiceInfo{
			ServiceId: message.ServiceId,
			Path:      message.Path,
			Port:      message.Port,
			IP:        message.Port,
		}

		if err := rt.balancer.AddService(newService); err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
			return
		}

		updatedServices := rt.registry.GetServices()

		listComponent := make([]tmpl.Service, 0)

		for _, updatedService := range updatedServices {
			listComponent = append(listComponent,
				tmpl.Service{
					Port:      updatedService.Port,
					Path:      updatedService.Path,
					ServiceId: updatedService.ServiceId,
					IP:        updatedService.IP,
					IsHealthy: updatedService.IsHealthy,
				})
		}

		buffer := new(bytes.Buffer)
		comp := tmpl.ServiceListComponent(listComponent)

		err = comp.Render(context.Background(), buffer)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
		}

		rt.hub.Broadcaster() <- buffer.Bytes()
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
			newSrvComp := tmpl.Service{
				Port:      val.Port,
				Path:      val.Path,
				IsHealthy: val.IsHealthy,
				IP:        val.IP,
				ServiceId: val.ServiceId,
			}
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
			slog.Error(fmt.Sprint("Error upgrading to WebSocket: \n", err))
			return
		}

		defer conn.Close()

		newClient := NewSocketClient(rt.hub, conn, WithWriteWaitTime(10*time.Second))
		rt.hub.Register() <- &newClient

		newClient.ReadPipe(rt.ctx)
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

// WithSecretKey is an option for seting the hashSecretKey of the
// Mux router. Whatever is passed in hashed then stored.
func WithSecretKey(key string) MuxRouterOpt {
	return func(router *MuxRouter) error {
		if len(strings.TrimSpace(key)) == 0 {
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

func WithHub(hub Hub) MuxRouterOpt {
	return func(mr *MuxRouter) error {
		mr.hub = hub
		return nil
	}
}

// MuxRouterOpt are option functions that setup the mux router struct.
type MuxRouterOpt func(*MuxRouter) error

// New MuxRouter instantiates a MuxRouter with all the necessary handleFuncs utilizing the
// service registry. The MuxRouter implements the Router interface for
func NewMuxRouter(balancer balancer.LoadBalancer, registry registry.Registry, ctx context.Context, opts ...MuxRouterOpt) (Router, error) {
	router := &MuxRouter{
		balancer: balancer,
		registry: registry,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		ctx: ctx,
		hub: NewInMemoryHub(),
	}

	for _, opt := range opts {
		err := opt(router)
		if err != nil {
			return nil, err
		}
	}

	go router.hub.Run(ctx)

	return router, nil
}
