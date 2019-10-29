package notify

import (
	"log"
	"net/http"
	"strconv"

	"github.com/leaftree/onoffice/models"
)

func CounterNotice(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("userid")
	cardDate := r.FormValue("card_date")

	if uid == "" || cardDate == "" {
		w.Write([]byte("invalid request"))
		return
	}

	id, _ := strconv.ParseUint(uid, 10, 64)
	err := models.CounterNotice(id, cardDate)
	if err != nil {
		log.Println(err)
		w.Write([]byte("request failed"))
		return
	}
	w.Write([]byte("notify canceled"))
}
