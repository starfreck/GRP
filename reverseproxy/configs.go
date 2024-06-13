package reverseproxy

import (
	"errors"
	"fmt"
	"os"

	"strconv"
)

var (
	DEFAULT_PORT = "80"
)

func GetUpstreamServers() *map[string]*Server {

	var configs []*Server

	for i := 1; ; i++ {

		host := os.Getenv("TARGET_URL_" + strconv.Itoa(i))

		if host == "" {
			break
		}

		testCookie, _ := strconv.ParseBool(os.Getenv("TEST_COOKIE_" + strconv.Itoa(i)))

		ssl, _ := strconv.ParseBool(os.Getenv("SSL_" + strconv.Itoa(i)))

		config := &Server{}

		config.Host = host
		config.TestCookie = testCookie
		config.SSL = ssl

		configs = append(configs, config)
	}

	serverMap := make(map[string]*Server)
	for _, config := range configs {
		serverMap[config.Host] = config
	}

	i := 1
	for _, server := range serverMap {
		fmt.Printf("Loading target server config %d\n", i)
		fmt.Printf("Target Server: %+v, SSL: %v, TestCookie: %v\n", server.Host, server.SSL, server.TestCookie)
		i++
		if i <= len(serverMap) {
			fmt.Println()
		}
	}

	fmt.Printf("\nTotal %d servers loaded...\n", len(serverMap))

	if len(serverMap) == 0 {
		fmt.Fprintf(os.Stderr, "Error: %v\n", "No target servers are configured in .env file")
		os.Exit(1)
	}

	return &serverMap
}

func GetSSLHeadersToDiscard() ([]string, error) {

	var headers []string
	for i := 1; ; i++ {

		header := os.Getenv("DISCARD_SSL_HEADER_" + strconv.Itoa(i))

		if header == "" {
			break
		}

		headers = append(headers, header)

	}

	if len(headers) == 0 {
		return headers, errors.New("no SSL headers were found in the environment file")
	}

	return headers, nil
}

func GetPort() string {

	port := os.Getenv("PORT")

	if len(port) == 0 {
		return DEFAULT_PORT
	} else {
		return port
	}
}

func GetName() string {

	name := os.Getenv("NAME")

	if len(name) == 0 {
		name = "Go"
	}

	return name
}
