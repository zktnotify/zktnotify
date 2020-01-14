package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	jsonresp "github.com/zktnotify/zktnotify/pkg/resp"
	"github.com/zktnotify/zktnotify/router/notify"
	"github.com/zktnotify/zktnotify/router/server"
	"github.com/zktnotify/zktnotify/router/usermng"
	"github.com/zktnotify/zktnotify/router/wxpusher"
)

func NewApiMux() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonresp.Respond(w, 404, nil, "迷路了吧")
	})

	regRouter(r)

	r.PathPrefix("/").
		Handler(http.StripPrefix("/", http.FileServer(http.Dir("dist"))))

	return r
}

func regRouter(r *mux.Router) {
	v1s := r.PathPrefix("/api/v1").Subrouter()

	v1s.HandleFunc("/status", server.Status).Methods("GET")
	v1s.HandleFunc("/shutdown", server.Shutdown).Methods("GET")

	v1s.HandleFunc("/counternotice", notify.CounterNotice).Methods("GET")
	v1s.HandleFunc("/user", usermng.AddUser).Methods("POST")
	v1s.HandleFunc("/user", usermng.GetUsers).Methods("GET")
	v1s.HandleFunc("/user/{jobId}", usermng.GetUser).Methods("GET")
	v1s.HandleFunc("/user/{jobId}", usermng.DeleteUser).Methods("DELETE")
	v1s.HandleFunc("/user/{jobId}", usermng.ChangeUserStatus).Methods("PUT")

	r.HandleFunc("/api/wxpusher/signup", wxpusher.Signup).Methods("POST")
	r.HandleFunc("/api/wxpusher/follow/callback", wxpusher.Follow).Methods("POST")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
