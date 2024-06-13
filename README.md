# 🚀 GRP (Go Reverse Proxy) 🔄

GRP is a lightweight and efficient reverse proxy written in Go. It's designed with a unique feature to bypass the [testcookie-nginx-module](https://github.com/kyprizel/testcookie-nginx-module), making it a versatile tool for various web environments.

### 🛠️ Environment Configuration

Configure the following options in your `.env` file at the root of your project and you should load `.env` in your `main()` function:

```env
# Proxy server configuration
NAME=BOSE                                     # Proxy server nickname (use as metadata)
PROXY_HOST_ADDRESS=proxy.server.com           # The URL of the proxy server
SSL=true                                      # If the proxy server is using SSL certificate
PORT=8080                                     # The port on which the proxy server will run
API_KEY=                                      # If the target server has some kind of authentication
USER_AGENT=                                   # Add a user agent yourself
                                              # It is important where TEST_COOKIE is set to true
# Target server(s) configurations
TARGET_URL_1=first.server.com                 # The URL of the upstream server
TEST_COOKIE_1=true                            # If the target server is using testcookie-nginx-module
SSL_1=true                                    # If the target server is using SSL certificate

TARGET_URL_2=second.server.com
TEST_COOKIE_2=true
SSL_2=true

#     .
#     .
#     . 
# TARGET_URL_N=nth.server.com
# TEST_COOKIE_N=true
# SSL_N=true

# If the target server does not have SSL enabled
# below SSL headers will be remove from the request
# you can add your own headers
DISCARD_SSL_HEADER_1=Strict-Transport-Security
DISCARD_SSL_HEADER_2=Public-Key-Pins
DISCARD_SSL_HEADER_3=Upgrade-Insecure-Requests
DISCARD_SSL_HEADER_4=Content-Security-Policy
DISCARD_SSL_HEADER_5=Access-Control-Allow-Origin
DISCARD_SSL_HEADER_6=WWW-Authenticate
DISCARD_SSL_HEADER_7=X-Content-Type-Options
DISCARD_SSL_HEADER_8=X-Frame-Options
DISCARD_SSL_HEADER_9=X-XSS-Protection
```
## ⚠️ Important Note 
When adding URLs, please omit the protocols (i.e., do not use http://example.com or https://example.com, but simply use `example.com`). The protocol will be automatically added based on the values of the SSL variables. If you do not specify SSL variables for a particular target server, the default value will be `‘false’`.

## 📖 Usage Examples

Please check out the [examples](./examples/) folder for more examples.

### Simple Load Balancer

```go
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
var proxyRequestMiddleware = func(s *reverseproxy.Server, proxyReq *http.Request, sourceReq *http.Request) {
	// Add additional headers to the proxy request
	proxyReq.Header.Set("Accept", "*/*")
	// Set Auth credentials if you are using any
	if len(os.Getenv("API_KEY")) > 0 {
		proxyReq.Header.Set("Authorization", os.Getenv("API_KEY"))
	}
}

```

## 🤝 Contribute

Contributions are always welcome from the community. Don’t hesitate to open an issue or submit a pull request!

## 📄 License

GRP is licensed under the [BSD-4-Clause](./LICENSE). See [LICENSE](./LICENSE) for more information.