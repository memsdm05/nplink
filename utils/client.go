package utils

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	http.Client
	Header http.Header
}

func NewClient() *Client {
	c := new(Client)
	c.Jar, _ = cookiejar.New(nil)
	c.Header = make(http.Header)
	c.Header.Set("user-agent", "nplink")
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

func (c *Client) DoJSON(req *http.Request, v interface{}) (*http.Response, error) {
	r, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if json.NewDecoder(r.Body).Decode(&v) != nil {
		return nil, err
	}
	defer r.Body.Close()

	return r, nil
}