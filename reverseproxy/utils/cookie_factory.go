package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type CookieFactory struct {
	Url        string
	Cookie     string
	RefererURL string
}

func (cf *CookieFactory) GetCookie(sourceReqHeaders http.Header) error {
	req, err := http.NewRequest("GET", cf.Url, nil)
	if err != nil {
		return err
	}

	req.Header = sourceReqHeaders

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	html := string(body)
	cf.extractCookieFromHTML(html)

	return nil
}

func (cf *CookieFactory) extractCookieFromHTML(html string) {
	re := regexp.MustCompile(`toNumbers\("([a-fA-F0-9]+)"\)`)
	matches := re.FindAllStringSubmatch(html, -1)

	a, _ := hex.DecodeString(matches[0][1])
	b, _ := hex.DecodeString(matches[1][1])
	c, _ := hex.DecodeString(matches[2][1])

	block, err := aes.NewCipher(a)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, b)
	mode.CryptBlocks(c, c)

	reLocation := regexp.MustCompile(`location.href="([^"]+)"`)
	locationMatches := reLocation.FindStringSubmatch(html)

	if len(locationMatches) > 1 {
		cf.Url = locationMatches[1]
		cf.RefererURL = locationMatches[1]
	}

	cf.Cookie = fmt.Sprintf("%x", c)
}
