package resp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func RenderJSON(w http.ResponseWriter, data JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(data)
	w.Write(bytes)
}

func RenderJSONOK(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	resp := JSONResponse{
		Status:  200,
		Message: "success",
		Data:    data,
	}
	bytes, _ := json.Marshal(resp)
	w.Write(bytes)
}

func Respond(w http.ResponseWriter, code int, data interface{}, msg ...interface{}) {
	rep := JSONResponse{
		Status:  code,
		Data:    data,
		Message: fmt.Sprint(msg),
	}
	if rep.Message == "" {
		rep.Message = "success"
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(rep)
	w.Write(bytes)
}

func RespondText(w http.ResponseWriter, data interface{}) {
	switch data.(type) {
	case string:
		w.Write([]byte(data.(string)))
	case []byte:
		w.Write(data.([]byte))
	}
}
