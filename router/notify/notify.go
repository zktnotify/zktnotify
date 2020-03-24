package notify

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/zktnotify/zktnotify/models"
)

var (
	tpl *template.Template
)

type CancelNotification struct {
	Text      string
	TextColor string
}

func init() {
	var err error

	tpl, err = template.New("CancelNotification").Parse(CancelNotificationTemplate)
	if err != nil {
		panic(err)
	}
}

func makeRealText(text CancelNotification) []byte {
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, text); err != nil {
		log.Println(err)
		return []byte("500 Server internal error")
	}
	return buf.Bytes()
}

func CounterNotice(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("userid")
	cardDate := r.FormValue("card_date")

	if uid == "" || cardDate == "" {
		w.Write(makeRealText(CancelNotification{
			Text:      "Invalid request",
			TextColor: "#FF0000",
		}))
		return
	}

	id, _ := strconv.ParseUint(uid, 10, 64)

	if err := models.CounterNotice(id, cardDate); err != nil {
		log.Println(err)
		w.Write(makeRealText(CancelNotification{
			Text:      "Notify cancel failed.",
			TextColor: "#FF0000",
		}))
		return
	}

	w.Write(makeRealText(CancelNotification{
		Text:      "Notify cancel successful!",
		TextColor: "#000000",
	}))
}
