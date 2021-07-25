package nightbot

import (
	"net/http"
	"net/url"
)

type nightbotService struct {
	client http.Client
	commandIds map[string]string
	session string
}

func (n *nightbotService) AuthURL() string {
	panic("implement me")
}

func (n *nightbotService) SetSession(session string) {
	panic("implement me")
}

func (n *nightbotService) CheckSession() bool {
	panic("implement me")
}

func (n *nightbotService) ResolveSession(vals url.Values) (string, error) {
	panic("implement me")
}

func (n nightbotService) SetCommand() {
	panic("implement me")
}

func (n nightbotService) DeleteCommand() {
	panic("implement me")
}
