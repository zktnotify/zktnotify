package zkt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	gourl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/opesun/goquery"
	"github.com/zktnotify/zktnotify/models"
	"github.com/zktnotify/zktnotify/pkg/xhttp"
)

var (
	LoginURL   string
	UserIDURL  string
	TimeTagURL string
)

func RegisterURL(base, login, userid, timetag string) (err error) {
	if base == "" {
		return errors.New("zkt server host is required")
	}
	LoginURL = base + login
	UserIDURL = base + userid
	TimeTagURL = base + timetag

	return nil
}

func Login(username string, userID uint64, password string) error {
	uparam := gourl.Values{}
	uparam.Add("username", username)
	uparam.Add("password", password)
	uparam.Add("login_type", "pwd")
	urldata := uparam.Encode()

	client := http.Client{}
	req, _ := http.NewRequest("POST", LoginURL, strings.NewReader(urldata))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if string(data) != "ok" {
		return errors.New(string(data))
	}

	if cookies := resp.Cookies(); len(cookies) > 0 {
		ck := http.Cookie{
			Name:  cookies[0].Name,
			Path:  cookies[0].Path,
			Value: cookies[0].Value,
		}
		CookieSet(username, userID, &ck)
	}
	return nil
}

func GetUserID(username string) (uid uint64, err error) {
	ck, ok := CookieGet(username, 0)
	if !ok {
		return 0, fmt.Errorf("cookie not found for:%s", username)
	}
	data, err := xhttp.GetWithCookie(UserIDURL, &ck, nil)
	if err != nil {
		return uid, err
	}

	nodes, err := goquery.Parse(bytes.NewBuffer(data))
	if err != nil {
		return uid, err
	}

	next := false
	nodes.Find("body input").Each(func(i int, s *goquery.Node) {
		for _, attr := range s.Attr {
			if attr.Val == "id_self_services" {
				next = true
			} else if next == true && attr.Key == "value" {
				uid, _ = strconv.ParseUint(attr.Val, 10, 64)
				break
			}
		}
	})

	if uid == 0 {
		return 0, errors.New("user id not found by url parse")
	}

	CookieUpdate(username, uid)

	return uid, nil
}

func GetTimeTag(uid uint64, start, end time.Time) (_ *models.TimeTagInfos, _ error) {
	var (
		infos     = models.TimeTagInfos{}
		uparam    = gourl.Values{}
		client    = http.Client{}
		cookie, _ = CookieGet("", uid)
	)

	uparam.Add("page", "1")
	uparam.Add("rp", "20")
	uparam.Add("isForce", "0")
	uparam.Add("UserIDs", fmt.Sprintf("%d", uid))
	uparam.Add("ComeTime", start.Format("2006-01-02"))
	uparam.Add("EndTime", end.Format("2006-01-02"))
	urldata := uparam.Encode()

	req, _ := http.NewRequest("POST", TimeTagURL, strings.NewReader(urldata))
	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &infos); err != nil {
		if strings.Contains(string(data), "!DOCTYPE HTML") {
			err = nil
		}
		return nil, err
	}
	return &infos, nil
}
