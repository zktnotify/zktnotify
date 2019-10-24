package zkt

import (
	"errors"
	"strings"
	"time"

	"github.com/imroc/req"
	"github.com/leaftree/onoffice/models"
)

func Login(url, username, password string) error {
	header := req.Header{}
	param := req.Param{
		"username":   username,
		"password":   password,
		"login_type": "pwd",
	}

	r, err := req.Post(url, header, param)
	if err != nil {
		return err
	}
	if r.String() == "ok" {
		return nil
	}
	return errors.New(r.String())
}

func GetTimeTag(url string, uid uint64, start, end time.Time) (_ *models.TimeTagInfos, _ error) {
	data := models.TimeTagInfos{}
	header := req.Header{}
	param := req.Param{
		"page":     1,
		"rp":       20,
		"UserIDs":  uid,
		"isForce":  0,
		"ComeTime": start.Format("2006-01-02"),
		"EndTime":  end.Format("2006-01-02"),
	}

	r, err := req.Post(url, header, param)
	if err != nil {
		return nil, err
	}

	if err := r.ToJSON(&data); err != nil {
		if strings.Contains(r.String(), "!DOCTYPE HTML") {
			err = nil
		}
		return nil, err
	}
	return &data, nil
}

func OvertimeReturn(over string) string {
	if over >= "20:00:00" && over < "21:00:00" {
		return "09:30:00"
	} else if over >= "21:00:00" && over < "22:00:00" {
		return "10:30:00"
	} else if over >= "22:00:00" {
		return "11:00:00"
	}
	return "09:15:00"
}
