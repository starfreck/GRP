package handlers

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/starfreck/grp/reverseproxy"
	"github.com/starfreck/grp/reverseproxy/middlewares"
)

type ReverseProxyHandler struct {
	upstreamServers    *map[string]*reverseproxy.Server
	proxyRequestRouter middlewares.ProxyRequestRouter
}

func NewReverseProxyHandler(upstreamServers *map[string]*reverseproxy.Server, proxyRequestRouter middlewares.ProxyRequestRouter, proxyRequestMiddleware middlewares.ProxyRequestMiddleware, proxyResponseMiddleware middlewares.ProxyResponseMiddleware) *ReverseProxyHandler {
	for host, server := range *upstreamServers {
		if proxyRequestMiddleware != nil {
			server.ProxyRequestMiddleware = proxyRequestMiddleware
		}
		if proxyResponseMiddleware != nil {
			server.ProxyResponseMiddleware = proxyResponseMiddleware
		}
		(*upstreamServers)[host] = server
	}
	return &ReverseProxyHandler{upstreamServers: upstreamServers, proxyRequestRouter: proxyRequestRouter}
}

func (rph *ReverseProxyHandler) HandleResource(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s: Received %s request at %s with %s", r.Method, r.Method, r.URL.Path, r.URL.Query())
	upstreamServer := rph.getUpstreamServer(*r)
	err := upstreamServer.RoundTrip(w, r, r.Method)
	// if there is any error on proxy server write to response
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: " + err.Error()))
		log.Printf("Error: %v\n", err.Error())
	}

}

// Get all upstream server hosts into a slice
func (rph ReverseProxyHandler) getUpstreamServerHosts() *[]string {

	hosts := make([]string, 0, len(*rph.upstreamServers))
	for k := range *rph.upstreamServers {
		hosts = append(hosts, k)
	}
	return &hosts
}

func (rph ReverseProxyHandler) getRandomUpstreamServerHost() string {

	hosts := *rph.getUpstreamServerHosts()

	// Create a new random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Pick a random key
	randomUpstreamServerHost := hosts[r.Intn(len(hosts))]

	return randomUpstreamServerHost
}

func (rph ReverseProxyHandler) getUpstreamServer(r http.Request) *reverseproxy.Server {

	// Get a random upstream server
	upstreamServerHost := rph.getRandomUpstreamServerHost()

	// Call the proxy server router if it's available
	if rph.proxyRequestRouter != nil {
		upstreamServerHost = rph.proxyRequestRouter(*rph.getUpstreamServerHosts(), r)
	}

	return (*rph.upstreamServers)[upstreamServerHost]
}
