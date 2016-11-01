package cores

import (
	"math/rand"
	"strings"
	"fmt"
	"net/http"
	"log"
)

var isFirst bool = true
const Fm = fmt.Sprintf

type Root struct{ Server string `xml:"server"` }

func RandInt(min, max int) int {
	return rand.Intn(max - min) + min
}

func Network(session *http.Client, url, method, query, referer string) (*http.Response, error) {
	var err error
	var req *http.Request
	switch method {
	case "GET":
		req, err = http.NewRequest("GET", url, nil)
		req.URL.RawQuery = query
	case "POST":
		req, err = http.NewRequest("POST", url, strings.NewReader(query))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36")
	if referer != "" {
		req.Header.Set("Referer", fmt.Sprintf("http://live.bilibili.com/%s", referer))
	}
	return session.Do(req)
}