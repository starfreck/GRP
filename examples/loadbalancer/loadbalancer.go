package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/starfreck/grp/reverseproxy"
	"github.com/starfreck/grp/reverseproxy/handlers"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Panicf("Environment file is not provided: %v", err)
	}

	// Create Proxy servers by reading the configuration
	upstreamServers := reverseproxy.GetUpstreamServers()
	// Create a proxy server handler
	reverseProxyHandler := handlers.NewReverseProxyHandler(upstreamServers, proxyRequestRouter, proxyRequestMiddleware, nil)

	router := http.NewServeMux()
	router.HandleFunc("/", reverseProxyHandler.HandleResource)

	fmt.Printf("🚀🔄 %v load balancer is starting on port %v...\n\n", reverseproxy.GetName(), reverseproxy.GetPort())
	log.Fatal(http.ListenAndServe(":"+reverseproxy.GetPort(), router))
}

func randomInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

// Select a target host based on the logic defined in this function
var proxyRequestRouter = func(upstreamHosts []string, r http.Request) string {
	randomIndex := randomInt(len(upstreamHosts))
	return upstreamHosts[randomIndex]
}

// Update request before forwarding to the target host
var proxyRequestMiddleware = func(s *reverseproxy.Server, sourceReq *http.Request) {
	// Add additional headers to the proxy request
	sourceReq.Header.Set("Accept", "*/*")
	// Set Auth credentials if you are using any
	if len(os.Getenv("API_KEY")) > 0 {
		sourceReq.Header.Set("Authorization", os.Getenv("API_KEY"))
	}
}
