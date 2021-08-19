package streamelements

import (
	"errors"
	"fmt"
	//"github.com/memsdm05/nplink/provider"
	"github.com/memsdm05/nplink/utils"
	"net/http"
	"net/url"
)

//func init() {
//	provider.Register(new(streamElementsService))
//}

const base =  "https://api.streamelements.com"

type streamElementsService struct {
	client *utils.Client
}

func (s *streamElementsService) Name() string {
	return "streamelements"
}

func (s *streamElementsService) Init() {
	s.client = utils.NewClient()
	s.client.CheckRedirect = utils.StopRedirect
	s.client.Header.Set("referer", "https://streamelements.com/")
}

func (s *streamElementsService) Session(session string) error {
	fmt.Println(session)
	panic("fuck me")
}

func (s *streamElementsService) URL() string {
	req, _ := http.NewRequest("GET", base + "/auth/twitch", nil)
	resp, _ := s.client.Do(req)
	return resp.Header.Get("location")
}

func (s *streamElementsService) ResolveSession(vals url.Values) (string, error) {
	seUrl := &url.URL{
		Scheme: "https",
		Host: "api.streamelements.com",
		Path: "/auth/twitch",
		RawQuery: vals.Encode(),
	}

	seUrl.Redacted()
	s.client.Get(seUrl.String())

	for _, cookie := range s.client.Jar.Cookies(seUrl) {
		if cookie.Name == "se-token" {
			fmt.Println(cookie.Value)
			return cookie.Value, nil
		}
	}
	return "", errors.New("cookie not found")
}

func (s *streamElementsService) SetCommand(name, msg string) error {
	panic("implement me")
}

func (s *streamElementsService) DeleteCommand(name string) error {
	panic("implement me")
}


