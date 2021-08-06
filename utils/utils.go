package utils

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Client struct {
	http.Client
	Header http.Header
}

func NewClient() *Client {
	c := new(Client)
	c.Jar, _ = cookiejar.New(nil)
	c.Header = make(http.Header)
	return c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	for k, v1 := range c.Header {
		for _, v2 := range v1 {
			req.Header.Add(k, v2)
		}
	}

	return c.Client.Do(req)
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

var StopRedirect = func(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

func TransposeValues(values url.Values, keys ...string) (ret url.Values) {
	for _, key := range keys { ret.Set(key, values.Get(key)) }
	return
}

func TransposeHeader(header http.Header, keys ...string) (ret http.Header) {
	for _, key := range keys { ret.Set(key, header.Get(key)) }
	return
}