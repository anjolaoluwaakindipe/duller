package routing

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

type Router interface {
	RegisterRoutes()
	GetRouter() http.Handler
}

type MuxRouter struct {
	router *mux.Router	
} 

func (mr *MuxRouter) RegisterRoutes() {
	mr.router.HandleFunc("/{path}",mr.GetPath(mr.ProxyRequest)).Methods("GET")
	mr.router.Use(mux.CORSMethodMiddleware(mr.router))
}

type MuxRoutes struct{

}
func (mr *MuxRouter) GetPath(proxyfunc func (string )(*httputil.ReverseProxy, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		fmt.Println(path)
		//  use path to get address from service discovery
		address := "asdf"

		proxy, err := proxyfunc(address)
		if err != nil {
			log.Printf("address of discovered service is invalid : %v", err)
			return 
		}

		proxy.ServeHTTP(w, r)

	}
}

func (mr *MuxRouter) GetAddress(address string) (string, error) {
	return "", nil 
}

func (mr *MuxRouter) ProxyRequest(targetUrl string) (*httputil.ReverseProxy, error){
	url , err := url.Parse(targetUrl)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

func (mr *MuxRouter) GetRouter() http.Handler {
	return mr.router
}

func InitMuxRouter() Router {
	return &MuxRouter{
		router: mux.NewRouter(),
	}
}