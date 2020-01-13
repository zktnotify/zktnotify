package wxpusher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zktnotify/zktnotify/models"
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
		w.Write([]byte("Error, TODO"))
		return
	}

	if err := json.Unmarshal(data, &reg); err != nil {
		w.Write([]byte("Error, TODO"))
		return
	}

	user := models.GetUserByJobId(reg.Account)
	if user != nil {
		w.Write([]byte("账号注册过了, TODO"))
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
		w.Write([]byte("账号注册失败了, TODO"))
		return
	}
	w.Write([]byte("账号成功了, TODO"))
}
