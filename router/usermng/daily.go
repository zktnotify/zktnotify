package usermng

import (
	"net/http"
	"strconv"

	"github.com/zktnotify/zktnotify/pkg/resp"
	"github.com/zktnotify/zktnotify/pkg/service"
)

func GetMonthDaily(w http.ResponseWriter, r *http.Request) {
	userid, _ := strconv.ParseUint(r.FormValue("userid"), 0, 64)
	month, _ := strconv.ParseInt(r.FormValue("month"), 0, 64)

	if userid == 0 {
		resp.RenderJSON(w, resp.InvalidRequest)
		return
	}

	data, err := service.RetrieveMonthDaily(userid, int(month))
	if err != nil {
		resp.RenderJSON(w, resp.ServiceError)
		return
	}
	resp.RespondText(w, data)
}
