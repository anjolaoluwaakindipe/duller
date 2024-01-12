package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

// Router requires some sort of implementation
type Router interface {
	RegisterRoutes()
	GetRouter() http.Handler
}

// MuxRouter this is a Gorilla Mux router implementation of the router needed for the gateway
type MuxRouter struct {
	router           *mux.Router
	discoveryAddress string
}

// RegisterRoutes registers all handlers needed for the gateway
func (mr *MuxRouter) RegisterRoutes() {
	mr.router.HandleFunc("/{path}", mr.GetPath(utils.ProxyRequest)).Methods("GET")
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

		address, status, getAddErr := mr.GetAddress(path)

		if getAddErr != nil {
			w.WriteHeader(status)
			response := GatewayErrorMessage{Message: getAddErr.Error(), Status: status}
			jsonResponse, _ := json.Marshal(&response)
			w.Write(jsonResponse)
			return
		}

		proxy, err := proxyfunc(address)
		if err != nil {
			log.Printf("address of discovered service is invalid : %v", err)
			return
		}

		proxy.ServeHTTP(w, r)
	}
}

// GetAddress messages a discovery server to get the appropriate service address
// for a requested path
func (mr *MuxRouter) GetAddress(path string) (string, int, error) {
	conn, err := net.Dial("tcp", mr.discoveryAddress)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	defer conn.Close()

	message := discovery.Message{
		Type: "getAddressMsg",
		Data: discovery.GetServiceMessage{Path: path},
	}

	jsonMessage, _ := json.Marshal(message)

	if _, err := conn.Write(jsonMessage); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error while sending path to registry")
	}

	decoder := json.NewDecoder(conn)

	var response discovery.RegistryResponse

	decoder.Decode(&response)

	if response.Code != 0 {
		return "", http.StatusNotFound, fmt.Errorf(response.Message)
	}

	var info struct {
		Address string `json:"address"`
	}

	if err := mapstructure.Decode(response.Data, &info); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("could not decode registry server response data")
	}

	return info.Address, http.StatusOK, nil
}

func (mr *MuxRouter) GetRouter() http.Handler {
	return mr.router
}

func InitMuxRouter(discoveryAddress string) Router {
	return &MuxRouter{
		router:           mux.NewRouter(),
		discoveryAddress: discoveryAddress,
	}
}
