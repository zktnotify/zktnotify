package wxpusher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/zktnotify/zktnotify/models"
	"github.com/zktnotify/zktnotify/pkg/config"
	"github.com/zktnotify/zktnotify/pkg/notify/wxpusher"
	"github.com/zktnotify/zktnotify/pkg/resp"
)

type WXPusherTokenCache struct {
	sync.Mutex
	token map[string]struct{}
}

var (
	wxpTokenCache = WXPusherTokenCache{token: map[string]struct{}{}}
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

	if models.IsTokenBind(reg.Token) {
		resp.RenderJSON(w, resp.JSONResponse{Status: 200, Message: "Token已经绑定过了"})
		return
	}

	wxpTokenCache.Lock()
	if _, ok := wxpTokenCache.token[reg.Token]; !ok {
		resp.RenderJSON(w, resp.JSONResponse{Status: 300, Message: "请先关注WXPusher公众号"})
		return
	}
	wxpTokenCache.Unlock()

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

type Follower struct {
	Action string `json:"action"`
	Data   struct {
		AppKey  string `json:"appKey"`
		AppName string `json:"appName"`
		Source  string `json:"source"`
		Time    int64  `json:"time"`
		UID     string `json:"uid"`
	} `json:"data"`
}

// Follow when user follow me, WXPusher will send me a follow message
func Follow(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	follow := Follower{}
	if err := json.Unmarshal(data, &follow); err != nil {
		log.Println("WXPusher关注回调失败：", err)
		return
	}

	const msg = `
# 消息通知
感谢关注，请点击[注册](http://zkt.fylos.cn:8000/#/auth/signup?uid={{.UID}})来完成注册

【提示】**点击查看原文** 后再进行注册
	`

	text := strings.Replace(msg, "{{.UID}}", follow.Data.UID, 1)

	notifier := wxpusher.New()
	notifier.SetAppToken(config.Config.XServer.NotificationServer.AppToken)

	if err = notifier.Notify(follow.Data.UID, text); err != nil {
		log.Println("通知用户完成关注失败：", err)
	}
	wxpTokenCache.Lock()
	wxpTokenCache.token[follow.Data.UID] = struct{}{}
	wxpTokenCache.Unlock()
}
