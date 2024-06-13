package reverseproxy

import (
	"errors"
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

	// TODO: Remove this outside of this struct or else it won't be shared
	CookieJar     map[string]string
	cookieFactory utils.CookieFactory

	ProxyRequestMiddleware  func(s *Server, req *http.Request, r *http.Request)
	ProxyResponseMiddleware func(s *Server, w http.ResponseWriter, resp *http.Response)
}

func (s *Server) init() {
	s.setProxyURL()
	s.CookieJar = make(map[string]string)
	s.cookieFactory = utils.CookieFactory{Url: s.Url, RefererURL: s.RefererURL, Cookie: ""}
}

func (s *Server) setProxyURL() {
	s.Scheme = "http"
	if s.SSL {
		s.Scheme = "https"
	}
	s.Url = fmt.Sprintf("%s://%s/", s.Scheme, s.Host)
}

func (s *Server) getUserAgent() (string, error) {

	header := os.Getenv("USER_AGENT")

	if len(header) == 0 {
		return header, errors.New("no user agent was found in the environment file, It will fail the request")
	}

	return header, nil
}

func (s *Server) RoundTrip(w http.ResponseWriter, r *http.Request, method string) error {

	s.init()

	log.Println("Target server URL: ", s.Url)

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
	if s.TestCookie && s.CookieJar[s.Host] == "" {

		// Set user-agent before creating the cookie request
		// Update the User-Agent header in original request
		userAgent, err := s.getUserAgent()
		if err != nil {
			log.Panicln(err)
		}
		r.Header.Set("User-Agent", userAgent)
		log.Printf("Auth cookie is not set for %s...\n", s.Url)
		// Clone the headers from the original request to get the cookie
		if err := s.cookieFactory.GetCookie(r.Header.Clone()); err != nil {
			return err
		} else {
			s.Url = s.cookieFactory.Url
			s.RefererURL = s.cookieFactory.RefererURL
			s.CookieJar[s.Host] = s.cookieFactory.Cookie
		}
	}

	if s.TestCookie {
		r.Header.Set("Cookie", fmt.Sprintf("__test=%s; expires=Thu, 31-Dec-37 23:55:55 GMT; path=/", s.CookieJar[s.Host]))
	}
	if s.RefererURL != "" {
		r.Header.Set("Referer", s.RefererURL)
	}

	// Set the host in request headers
	r.Header.Set("Host", s.Host)

	// Create a new proxy request
	req, err := http.NewRequest(method, s.Url, r.Body)
	if err != nil {
		return err
	}

	// Copy headers and query params from the source request to proxy request
	utils.CopySourceAttributes(req, r)

	// Call user callback function to modify the proxy request
	if s.ProxyRequestMiddleware != nil {
		s.ProxyRequestMiddleware(s, req, r)
	}

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
