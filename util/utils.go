package util

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

const credFileName = "cred.txt"

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
			req.Header.Set(k, v2)
		}
	}

	return c.Client.Do(req)
}

func Must(err error) {
	if err != nil {
		FatalError(err)
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

func FatalError(err error)  {
	fmt.Println("An error has occurred:\n")
	color.Red("%v", err)
	fmt.Println("\n\n< Press Enter to Exit >a")
	//fmt.Printf("%s\n\n< press enter to close >", err)
	fmt.Scanln()
	os.Exit(727) // I am very funny
}

func GetCred(prov string) (string, time.Time, error) {
	f, err := os.OpenFile(credFileName, os.O_RDONLY, 0o667)

	if err != nil {
		return "", time.Time{}, err
	}

	content, _ := io.ReadAll(f)

	l := strings.Split(string(content), "\n")
	t, err := time.Parse(time.RFC3339, l[2])

	if err != nil {
		return "", time.Time{}, err
	}

	if l[0] != prov {
		return "", time.Time{}, errors.New("incorrect provider")
	}

	return l[1], t, nil
}

func SetCred(prov, cred string) error {
	f, err := os.OpenFile(credFileName, os.O_WRONLY|os.O_CREATE, 0o667)
	defer f.Close()

	if err != nil {
		return err
	}


	f.WriteString(prov + "\n")
	f.WriteString(cred + "\n")
	f.WriteString(time.Now().Format(time.RFC3339))

	return nil
}