package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"strconv"

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
	reverseProxyHandler := handlers.NewReverseProxyHandler(upstreamServers, proxyRequestRouter, proxyRequestMiddleware, proxyResponseMiddleware)

	router := http.NewServeMux()
	router.HandleFunc("/", reverseProxyHandler.HandleResource)

	fmt.Printf("🚀🔄 %v reverse proxy server is starting on port %v...\n\n", reverseproxy.GetName(), reverseproxy.GetPort())
	log.Fatal(http.ListenAndServe(":"+reverseproxy.GetPort(), router))
}

// Select a target host based on the logic defined in this function
var proxyRequestRouter = func(upstreamHosts []string, r http.Request) string {
	// Implement logic to select a host based on the (read only) request
	// Parse query parameters from the URL
	params := r.URL.Query()
	// Access specific query parameter
	fileType := params.Get("filetype")
	// you can do  many things here
	// i.e. load balancing, db/cache queries etc and redirect requests to the specific target server
	if len(fileType) > 0 {
		if fileType == "file" {
			return upstreamHosts[1]
		} else if fileType == "image" {
			return upstreamHosts[2]
		}
	}
	return upstreamHosts[0]
}

// Update request before forwarding to the target host
var proxyRequestMiddleware = func(s *reverseproxy.Server, sourceReq *http.Request) {
	// Add additional headers to the proxy request
	sourceReq.Header.Set("Accept", "*/*")
	// Set Auth credentials
	if len(os.Getenv("API_KEY")) > 0 {
		sourceReq.Header.Set("Authorization", os.Getenv("API_KEY"))
	}
}

// Update the response before forwarding to source
var proxyResponseMiddleware = func(s *reverseproxy.Server, w http.ResponseWriter, resp *http.Response) {
	// Check the status code and content type
	// check if the received response is json and if it contains the target server URL then replace it with proxy server's URL
	if resp.StatusCode == http.StatusOK && resp.Header.Get("Content-Type") == "application/json" {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		// Close the response body since it will be empty one it is being read
		err = resp.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
		// Modify the host scheme in any incoming messages where the type is application/json
		ssl, err := strconv.ParseBool(os.Getenv("SSL"))
		if err != nil {
			log.Fatalln(err)
		}
		// log.Println("Response Before:", string(b))
		// Make sure to replace scheme properly
		proxyServerScheme := "http"
		if ssl {
			proxyServerScheme = "https"
		}
		b = bytes.Replace(b, []byte(s.Scheme), []byte(proxyServerScheme), -1)
		// Modify the host address in any incoming messages where the type is application/json
		if len(os.Getenv("PROXY_HOST_ADDRESS")) > 0 {
			b = bytes.Replace(b, []byte(s.Host), []byte(os.Getenv("PROXY_HOST_ADDRESS")), -1)
		}
		// log.Println("Response After:", string(b))
		body := io.NopCloser(bytes.NewReader(b))
		resp.Body = body
		resp.ContentLength = int64(len(b))
		// Update the content length since the response body is changed
		resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	}
}
