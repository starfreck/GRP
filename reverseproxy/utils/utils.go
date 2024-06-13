package utils

import (
	"io"
	"log"
	"net/http"
)

func CopySourceAttributes(proxyReq *http.Request, sourceReq *http.Request) {
	// Copy headers
	for name, headers := range sourceReq.Header {
		for _, h := range headers {
			proxyReq.Header.Set(name, h)
		}
	}
	// Copy query parameters
	q := proxyReq.URL.Query()
	for key, values := range sourceReq.URL.Query() {
		for _, value := range values {
			q.Set(key, value)
		}
	}
	proxyReq.URL.RawQuery = q.Encode()
}

func CopyProxyHeaders(sourceResponseWriter http.ResponseWriter, proxyResp *http.Response) {
	sourceResponseWriter.WriteHeader(proxyResp.StatusCode)
	for name, values := range proxyResp.Header {
		for _, value := range values {
			sourceResponseWriter.Header().Set(name, value)
		}
	}
}

func CopyProxyResponse(sourceResponseWriter http.ResponseWriter, proxyResp *http.Response) {
	if _, err := io.Copy(sourceResponseWriter, proxyResp.Body); err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}
