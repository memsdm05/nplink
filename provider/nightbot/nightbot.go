package nightbot

import (
	"github.com/memsdm05/nplink/provider"
	"github.com/memsdm05/nplink/util"
	"net/url"
)

func init() {
	provider.Register(new(nightbotService))
}

type nightbotService struct {
	client *util.Client
	commandIds map[string]string
}

func (n nightbotService) Name() string {
	return "nightbot"
}

func (n *nightbotService) Init() {
	n.client = util.NewClient()
	n.client.Header.Set("user-agent", "nightbot @ nplink (github.com/memsdm05/nplink)")
}

func (n *nightbotService) Session(session string) error {
	n.client.Header.Set("authorization", "Session " + session)
	panic("implement me")
}

func (n nightbotService) URL() string {
	panic("implement me")
}

func (n nightbotService) ResolveSession(vals url.Values) (string, error) {
	panic("implement me")
}

func (n nightbotService) SetCommand(name, msg string) error {
	panic("implement me")
}

func (n nightbotService) DeleteCommand(name string) error {
	panic("implement me")
}


