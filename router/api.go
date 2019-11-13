package router

import (
	"github.com/leaftree/ctnotify/router/usermng"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	jsonresp "github.com/leaftree/ctnotify/pkg/resp"
	"github.com/leaftree/ctnotify/router/notify"
	"github.com/leaftree/ctnotify/router/server"
)

func NewApiMux() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonresp.Respond(w, 404, nil, "迷路了吧")
	})

	regRouter(r)
	return r
}

func regRouter(r *mux.Router) {
	v1s := r.PathPrefix("/api/v1").Subrouter()

	v1s.HandleFunc("/status", server.Status).Methods("GET")
	v1s.HandleFunc("/shutdown", server.Shutdown).Methods("GET")

	v1s.HandleFunc("/counternotice", notify.CounterNotice).Methods("GET")
	v1s.HandleFunc("/usermng/user", usermng.AddUser).Methods("POST")
	v1s.HandleFunc("/usermng/user/{jobId}", usermng.GetUser).Methods("GET")
	v1s.HandleFunc("/usermng/user/{jobId}", usermng.DeleteUser).Methods("DELETE")
	v1s.HandleFunc("/usermng/user/{jobId}", usermng.ChangeUserStatus).Methods("PUT")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
