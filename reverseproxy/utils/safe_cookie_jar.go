package utils

import "sync"

type SafeCookieJar struct {
	sync.RWMutex
	cookieJar       map[string]string
	urlCache        map[string]string
	refererURLCache map[string]string
}

func NewSafeCookieJar() *SafeCookieJar {
	return &SafeCookieJar{cookieJar: make(map[string]string), urlCache: make(map[string]string), refererURLCache: make(map[string]string)}
}

func (sc *SafeCookieJar) SetCookie(key string, value string) {
	sc.Lock()
	defer sc.Unlock()
	sc.cookieJar[key] = value
}

func (sc *SafeCookieJar) GetCookie(key string) string {
	sc.RLock()
	defer sc.RUnlock()
	value, ok := sc.cookieJar[key]
	if !ok {
		return ""
	}
	return value
}

func (sc *SafeCookieJar) SetUrl(key string, value string) {
	sc.Lock()
	defer sc.Unlock()
	sc.urlCache[key] = value
}

func (sc *SafeCookieJar) GetUrl(key string) string {
	sc.RLock()
	defer sc.RUnlock()
	value, ok := sc.urlCache[key]
	if !ok {
		return ""
	}
	return value
}

func (sc *SafeCookieJar) SetRefererUrl(key string, value string) {
	sc.Lock()
	defer sc.Unlock()
	sc.refererURLCache[key] = value
}

func (sc *SafeCookieJar) GetRefererUrl(key string) string {
	sc.RLock()
	defer sc.RUnlock()
	value, ok := sc.refererURLCache[key]
	if !ok {
		return ""
	}
	return value
}
