package shorturl

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/leaftree/ctnotify/pkg/config"
)

func ShortURL(urlSuffix string, data map[string]interface{}) string {
	key := config.Config.XServer.ShortURL.Server.AppKey
	addr := config.Config.XServer.ShortURL.Server.ApiAddr

	raw := url.Values{}
	for k, v := range data {
		raw.Set(k, fmt.Sprint(v))
	}
	rawurl := config.Config.XServer.ShortURL.PrefixURL + urlSuffix + "?" + raw.Encode()

	raw = url.Values{}
	raw.Set("key", key)
	raw.Set("url", rawurl)
	//raw.Set("format", "json")
	raw.Set("expireDate", "2030-03-31")

	rep, err := http.Get(addr + "?" + raw.Encode())
	if err != nil {
		log.Println(err)
		return ""
	}
	defer rep.Body.Close()

	val, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	if !strings.HasPrefix(string(val), "http://") {
		return ""
	}

	return string(val)
}
