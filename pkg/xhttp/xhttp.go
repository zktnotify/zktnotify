package xhttp

import (
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

func Get(pathname string) ([]byte, error) {
	rsp, err := http.Get(pathname)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	return ioutil.ReadAll(rsp.Body)
}
