package routing

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

type Router interface {
	RegisterRoutes()
	GetRouter() http.Handler
}

type MuxRouter struct {
	router           *mux.Router
	discoveryAddress string
}

func (mr *MuxRouter) RegisterRoutes() {
	mr.router.HandleFunc("/{path}", mr.GetPath(mr.ProxyRequest)).Methods("GET")
	mr.router.Use(mux.CORSMethodMiddleware(mr.router))
}

type MuxRoutes struct {
}

func (mr *MuxRouter) GetPath(proxyfunc func(string) (*httputil.ReverseProxy, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		//  edit path
		path := vars["path"]
		if len(path) < 1 {
			response := GatewayErrorMessage{Message: "No path specified", Status: http.StatusBadRequest}
			jsonResponse, _ := json.Marshal(&response)
			w.Write(jsonResponse)
			return
		}

		if path[0] != '/' {
			path = "/" + path
		}

		for path[len(path)-1] == '/' {
			path = path[0 : len(path)-1]
		}

		//  use path to get address from service discovery

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

func (mr *MuxRouter) GetAddress(path string) (string, int, error) {
	conn, err := net.Dial("tcp", mr.discoveryAddress)

	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	defer conn.Close()

	message := discovery.Message{
		Type: "getAddressMsg",
		Data: discovery.GetAddressMessage{Path: path},
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

func (mr *MuxRouter) ProxyRequest(targetUrl string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetUrl)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
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
