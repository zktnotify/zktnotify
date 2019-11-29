package xhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func PostForm(pathname string, data url.Values) ([]byte, error) {
	rsp, err := http.PostForm(pathname, data)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Get(pathname string, header ...map[string]interface{}) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", pathname, nil)

	for _, h := range header {
		for k, v := range h {
			req.Header.Add(k, fmt.Sprint(v))
		}
	}
	rep, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	return ioutil.ReadAll(rep.Body)
}
