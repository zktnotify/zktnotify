package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leaftree/onoffice/router/notify"
)

func NewApiMux() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("迷路了吧"))
	})

	regRouter(r)
	return r
}

func regRouter(r *mux.Router) {
	v1s := r.PathPrefix("/api/v1").Subrouter()
	v1s.HandleFunc("/counternotice", notify.CounterNotice).Methods("GET")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
