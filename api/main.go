package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middlewareHandler struct {
	r *httprouter.Router
}

func NewMiddlewareHandler(r *httprouter.Router) http.Handler {
	m := middlewareHandler{}
	m.r = r
	return m
}

// kidnap the request to do something
func (m middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)
	router.POST("/user/:user_name", Login)

	return router
}

func main() {
	r := RegisterHandlers()
	http.ListenAndServe(":8000", r)
}
