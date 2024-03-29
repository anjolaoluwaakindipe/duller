package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/anjolaoluwaakindipe/duller/internal/utils"
	"github.com/gorilla/mux"
)

// Router requires some sort of implementation
type Router interface {
	RegisterRoutes()
	GetRouter() http.Handler
}

// MuxRouter this is a Gorilla Mux router implementation of the router needed for the gateway
type MuxRouter struct {
	router        *mux.Router
	discoveryHost string
	discoveryPort string
	discoveryPath string
}

// RegisterRoutes registers all handlers needed for the gateway
func (mr *MuxRouter) RegisterRoutes() {
	mr.router.HandleFunc("/{path}", mr.GetPath(utils.ProxyRequest))
	mr.router.Use(mux.CORSMethodMiddleware(mr.router))
}

// GetPath takes in a path variable from the gateway url and proxies the request
// to its associated service. It makes a request to the discovery server to fetch the right
// address to proxy the request to
func (mr *MuxRouter) GetPath(proxyfunc func(string) (*httputil.ReverseProxy, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		//  edit path
		path, ok := vars["path"]
		response := GatewayErrorMessage{}

		if !ok {
			response.Message = "No path specified"
			response.Status = http.StatusBadRequest
			jsonResponse, _ := json.Marshal(&response)
			w.WriteHeader(response.Status)
			w.Write(jsonResponse)
			return
		}

		utils.MakeUrlPathValid(&path)

		r.URL.Path = mr.discoveryPath + r.URL.Path
		proxy, err := proxyfunc("http://" + mr.discoveryHost + ":" + mr.discoveryPort)
		if err != nil {
			log.Printf("address of discovered service is invalid : %v", err)
			return
		}

		proxy.ServeHTTP(w, r)
	}
}

func (mr *MuxRouter) GetRouter() http.Handler {
	return mr.router
}

type MuxRouterOpts func(*MuxRouter)

func WithDiscoveryHost(discoveryAddress string) MuxRouterOpts {
	return func(mr *MuxRouter) {
		mr.discoveryHost = discoveryAddress
	}
}

func WithDiscoveryPort(discoveryPort string) MuxRouterOpts {
	return func(mr *MuxRouter) {
		mr.discoveryPort = discoveryPort
	}
}

func WithDiscoveryPath(discoveryPath string) MuxRouterOpts {
	return func(mr *MuxRouter) {
		mr.discoveryPath = discoveryPath
	}
}

func InitMuxRouter(opts ...MuxRouterOpts) Router {
	mr := &MuxRouter{
		router:        mux.NewRouter(),
		discoveryPort: "9876",
		discoveryHost: "localhost",
		discoveryPath: "/get-service/",
	}

	for _, opt := range opts {
		opt(mr)
	}

	return mr
}
