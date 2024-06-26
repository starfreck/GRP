package reverseproxy

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/starfreck/grp/reverseproxy/utils"
)

type Server struct {
	Host       string `json:"url"`
	TestCookie bool   `json:"test_cookie"`
	SSL        bool   `json:"ssl"`

	Url        string
	Scheme     string
	RefererURL string

	cookieFactory utils.CookieFactory

	ProxyRequestMiddleware  func(*Server, *http.Request)
	ProxyResponseMiddleware func(*Server, http.ResponseWriter, *http.Response)
}

func (s *Server) init() {
	s.setProxyURL()
	s.cookieFactory = utils.CookieFactory{Url: s.Url, RefererURL: s.RefererURL, Cookie: ""}
}

func (s *Server) setProxyURL() {
	s.Scheme = "http"
	if s.SSL {
		s.Scheme = "https"
	}
	s.Url = fmt.Sprintf("%s://%s/", s.Scheme, s.Host)
}

func (s *Server) getUserAgent() string {

	header := os.Getenv("USER_AGENT")

	if len(header) == 0 {
		log.Println("no user agent was found in the environment file, it could fail the request")
	}

	return header
}

func (s *Server) RoundTrip(w http.ResponseWriter, r *http.Request, method string, cookieJar *utils.SafeCookieJar) error {

	s.init()

	log.Println("Target server URL: ", s.Url)

	// Call user callback function to modify the proxy request
	if s.ProxyRequestMiddleware != nil {
		s.ProxyRequestMiddleware(s, r)
	}

	// Important Note:
	// If the target server does not support SSL, remove all SSL related headers
	// SSL headers can be configured in the environment file
	if !s.SSL {
		SSLHeaders, err := GetSSLHeadersToDiscard()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			for _, sslHeader := range SSLHeaders {
				r.Header.Del(sslHeader)
			}
		}
	}

	// Important Note:
	// If the target server has testcookie-nginx-module enabled set __test cookie
	// A fixed User-Agent is required and configured in the config file
	if s.TestCookie && (*cookieJar).GetCookie(s.Host) == "" {

		// Set User-Agent before creating the cookie request
		// Update the User-Agent header in original request
		r.Header.Set("User-Agent", s.getUserAgent())
		log.Printf("Auth cookie is not set for %s...\n", s.Url)
		if err := s.cookieFactory.GetCookie(s.getUserAgent()); err != nil {
			return err
		} else {
			s.Url = s.cookieFactory.Url
			s.RefererURL = s.cookieFactory.RefererURL
			// Cache cookie, url and referer url in cookie jar
			(*cookieJar).SetCookie(s.Host, s.cookieFactory.Cookie)
			(*cookieJar).SetUrl(s.Host, s.cookieFactory.Url)
			(*cookieJar).SetRefererUrl(s.Host, s.cookieFactory.RefererURL)

		}
	}

	// Get cookie, url and referer url from the cookie jar
	if s.TestCookie {
		r.Header.Set("User-Agent", s.getUserAgent())
		r.Header.Set("Cookie", fmt.Sprintf("__test=%s; expires=Thu, 31-Dec-37 23:55:55 GMT; path=/", (*cookieJar).GetCookie(s.Host)))
		s.Url = (*cookieJar).GetUrl(s.Host)
		s.RefererURL = (*cookieJar).GetRefererUrl(s.Host)
	}

	// Set the host in request headers
	r.Header.Set("Host", s.Host)
	r.Header.Set("Accept-Encoding", "gzip, deflate")
	if s.RefererURL != "" {
		r.Header.Set("Referer", s.RefererURL)
	}

	// Create a new proxy request
	req, err := http.NewRequest(method, s.Url, r.Body)
	if err != nil {
		return err
	}

	// Copy headers and query params from the source request to proxy request
	utils.CopySourceAttributes(req, r)

	// Make a proxy request call
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy headers from proxy response to source response
	utils.CopyProxyHeaders(w, resp)

	// Call user callback function to modify the source response
	if s.ProxyResponseMiddleware != nil {
		s.ProxyResponseMiddleware(s, w, resp)
	}

	// Send the response to source
	utils.CopyProxyResponse(w, resp)

	return nil
}
