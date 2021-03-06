package zkt

import (
	"net/http"
	"sync"
)

var (
	users    = map[uint64]string{}
	ncookies = map[string]*http.Cookie{}
	ucookies = map[uint64]*http.Cookie{}
	mux      = sync.Mutex{}
)

func CookieSet(jid string, uid uint64, ck *http.Cookie) {
	mux.Lock()
	defer mux.Unlock()

	if jid != "" {
		ncookies[jid] = ck
	}
	if uid != 0 {
		ucookies[uid] = ck
	}
	if uid != 0 && jid != "" {
		users[uid] = jid
	}
}

func CookieGet(name string, id uint64) (http.Cookie, bool) {
	mux.Lock()
	defer mux.Unlock()

	if id != 0 {
		ck, ok := ucookies[id]
		return *ck, ok
	}
	if name != "" {
		ck, ok := ncookies[name]
		return *ck, ok
	}
	return http.Cookie{}, false
}

func CookieUpdate(name string, id uint64) {
	mux.Lock()
	defer mux.Unlock()

	if ck, ok := ncookies[name]; id != 0 && ok {
		ucookies[id] = ck
		users[id] = name
	}
}

func CookieRemove(id uint64) {
	mux.Lock()
	defer mux.Unlock()

	if name, ok := users[id]; ok {
		delete(users, id)
		delete(ucookies, id)
		delete(ncookies, name)
	}
}

func HasCookie(jobID string, userID uint64) bool {
	mux.Lock()
	defer mux.Unlock()

	if userID != 0 {
		if _, ok := ucookies[userID]; ok {
			return true
		}
		return false
	}
	if jobID != "" {
		if _, ok := ncookies[jobID]; ok {
			return true
		}
	}
	return false
}
