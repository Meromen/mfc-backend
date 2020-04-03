package controller

import (
	"fmt"
	"github.com/Meromen/mfc-backend/db"
	"github.com/Meromen/mfc-backend/logger"
	"github.com/carlescere/scheduler"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"runtime"
)

const (
	API_URL string = "/api"
)

type response struct {
	Code   int
	Status string
	Body   interface{}
}

type controller struct {
	mfcStorage db.Storage
	logger     logger.Logger
}

type Controller interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func (c controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := c.NewRouter()
	router.ServeHTTP(w, r)
}

func (c controller) NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc(
		fmt.Sprintf("%s/mfcs", API_URL),
		headerMiddleware(c.GetMfcs, c.logger)).
		Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./client/")))
	return router
}

func NewController(storage db.Storage, logger logger.Logger) Controller {
	c := controller{storage, logger}
	scheduler.Every(1).Minutes().Run(c.UpdateMfcs)

	return &c
}

func headerMiddleware(handler http.HandlerFunc, logger logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		logger.Printf("Handler function called: %v", name)
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}
}
