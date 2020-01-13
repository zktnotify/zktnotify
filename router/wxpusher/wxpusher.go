package wxpusher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zktnotify/zktnotify/models"
	"github.com/zktnotify/zktnotify/pkg/resp"
)

type register struct {
	Name     string `json:"name"`
	Token    string `json:"token"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Account  string `json:"account"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	reg := register{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.RenderJSON(w, resp.JSONResponse{Status: 500, Message: "服务器端故障"})
		return
	}

	if err := json.Unmarshal(data, &reg); err != nil {
		resp.RenderJSON(w, resp.JSONResponse{Status: 300, Message: "请求参数异常"})
		return
	}

	user := models.GetUserByJobId(reg.Account)
	if user != nil {
		resp.RenderJSON(w, resp.JSONResponse{Status: 200, Message: "账号注册过了"})
		return
	}

	user = &models.User{
		Name:          reg.Name,
		JobID:         reg.Account,
		Password:      reg.Password,
		NotifyToken:   reg.Token,
		NotifyType:    2,
		NotifyAccount: reg.Mobile,
	}
	if err := models.SaveUser(user); err != nil {
		resp.RenderJSON(w, resp.JSONResponse{Status: 500, Message: "账号注册失败了:" + err.Error()})
		return
	}
	resp.RenderJSON(w, resp.JSONResponse{Status: 200, Message: "账号成功了"})
}
