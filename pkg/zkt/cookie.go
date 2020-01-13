package zkt

import (
	"net/http"
	"sync"
)

var (
	users   = map[string]uint64{}
	cookies = map[uint64]http.Cookie{}
	mux     = sync.Mutex{}
)

func CookieSet(jid string, uid uint64, ck http.Cookie) {
	mux.Lock()
	defer mux.Unlock()

	users[jid] = uid
	cookies[uid] = ck
}

func CookieGet(name string, id uint64) (http.Cookie, bool) {
	mux.Lock()
	defer mux.Unlock()

	if id != 0 {
		ck, ok := cookies[id]
		return ck, ok
	}
	if name != "" {
		if uid, ok := users[name]; ok {
			ck, ok := cookies[uid]
			return ck, ok
		}
	}
	return http.Cookie{}, false
}

func HasCookie(jobID string, userID uint64) bool {
	mux.Lock()
	defer mux.Unlock()

	if userID != 0 {
		if _, ok := cookies[userID]; ok {
			return true
		}
		return false
	}
	if jobID != "" {
		if uid, ok := users[jobID]; ok {
			if _, ok := cookies[uid]; ok {
				return true
			}
		}
	}
	return false
}
