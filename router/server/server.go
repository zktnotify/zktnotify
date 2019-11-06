package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/leaftree/ctnotify/pkg/config"
	jsonresp "github.com/leaftree/ctnotify/pkg/resp"
)

func Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonresp.RenderJSON(w, jsonresp.JSONResponse{
		Status:  0,
		Message: "server is running",
	})
}

func Shutdown(w http.ResponseWriter, r *http.Request) {
	jsonresp.RenderJSON(w, jsonresp.JSONResponse{
		Status:  0,
		Message: "server shutting down",
	})

	go func() {
		time.Sleep(time.Millisecond * 500)
		log.Println("server shutting down")
		os.Remove(config.Config.XServer.File.Pid)
		os.Exit(0)
	}()
}
