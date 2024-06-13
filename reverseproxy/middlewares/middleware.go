package middlewares

import (
	"net/http"

	"github.com/starfreck/grp/reverseproxy"
)

type ProxyRequestRouter func([]string, http.Request) string
type ProxyRequestMiddleware func(*reverseproxy.Server, *http.Request, *http.Request)
type ProxyResponseMiddleware func(*reverseproxy.Server, http.ResponseWriter, *http.Response)
