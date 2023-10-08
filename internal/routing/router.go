package routing

import (
	"fmt"
	"net/http"

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
	mr.router.HandleFunc("/{path}",mr.GetPath()).Methods("GET")
	mr.router.Use(mux.CORSMethodMiddleware(mr.router))
}

func (mr *MuxRouter) GetPath() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		fmt.Println(path)
	}
}

func (mr *MuxRouter) GetRouter() http.Handler {
	return mr.router
}

func InitMuxRouter() Router {
	return &MuxRouter{
		router: mux.NewRouter(),
	}
}