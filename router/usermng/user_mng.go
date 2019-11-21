package usermng

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zktnotify/zktnotify/pkg/resp"
	"github.com/zktnotify/zktnotify/pkg/service"
	"github.com/zktnotify/zktnotify/pkg/validate"
	"github.com/zktnotify/zktnotify/viewmodel"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("param read error"))
		return
	}
	user := new(viewmodel.User)
	if err = json.Unmarshal(body, user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("param error:%s\n", string(body))))
		return
	}
	if err := validate.VerifyStruct(user); err != nil {
		resp.RenderJSON(w, resp.JSONResponse{
			Status:  400,
			Message: "param check failed",
			Data:    validate.HandleVerifyErrorResult(err.(validator.ValidationErrors)),
		})
		return
	}
	userMng := service.GetUserManager()
	if err := userMng.AddUser(user); err != nil {
		resp.RenderJSON(w, resp.JSONResponse{
			Status:  500,
			Message: err.Error(),
		})
		return
	}
	resp.RenderJSON(w, resp.JSONResponse{
		Status:  200,
		Message: "OK",
	})
	return
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_jobId := vars["jobId"]
	if "" == strings.TrimSpace(_jobId) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the jobId is empty"))
		return
	}
	jobId, err := strconv.ParseUint(_jobId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("jobId error:%s", _jobId)))
		return
	}
	userMng := service.GetUserManager()
	user, err := userMng.GetUser(jobId)
	if err != nil {
		resp.RenderJSON(w, resp.JSONResponse{
			Status:  500,
			Message: err.Error(),
		})
		return
	}
	resp.RenderJSON(w, resp.JSONResponse{
		Status:  200,
		Message: "OK",
		Data:    user,
	})
	return
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_jobId := vars["jobId"]
	if "" == strings.TrimSpace(_jobId) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the jobId is empty"))
		return
	}
	jobId, err := strconv.ParseUint(_jobId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("jobId error:%s", _jobId)))
		return
	}
	userMng := service.GetUserManager()
	if err := userMng.DeleteUser(jobId); err != nil {
		resp.RenderJSON(w, resp.JSONResponse{
			Status:  500,
			Message: err.Error(),
		})
		return
	}
	resp.RenderJSON(w, resp.JSONResponse{
		Status:  200,
		Message: "OK",
	})
	return
}

func ChangeUserStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_jobId := vars["jobId"]
	if "" == strings.TrimSpace(_jobId) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the jobId is empty"))
		return
	}
	jobId, err := strconv.ParseUint(_jobId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("jobId error:%s", _jobId)))
		return
	}
	userMng := service.GetUserManager()
	if err := userMng.ChangeUserStatus(jobId); err != nil {
		resp.RenderJSON(w, resp.JSONResponse{
			Status:  500,
			Message: err.Error(),
		})
		return
	}
	resp.RenderJSON(w, resp.JSONResponse{
		Status:  200,
		Message: "OK",
	})
	return
}
