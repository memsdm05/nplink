package nightbot

import (
	"bytes"
	"encoding/json"
	"github.com/memsdm05/nplink/provider"
	"github.com/memsdm05/nplink/utils"
	"net/http"
	"net/url"
)

func init() {
	provider.Register(new(nightbotService))
}

const base = "https://api.nightbot.tv"

type nightbotService struct {
	client *utils.Client
	commandIds map[string]string
}

func (n nightbotService) Name() string {
	return "nightbot"
}

func (n *nightbotService) Init() {
	n.client = utils.NewClient()
	n.client.Header.Set("referer", "https://nightbot.tv/")
	n.client.Header.Set("origin", "https://nightbot.tv")

}

func (n *nightbotService) Session(session string) error {
	n.client.Header.Set("authorization", "Session " + session)
	req, _ := http.NewRequest("GET", base + "/1/commands", nil)
	resp, _ := n.client.Do(req)
	if resp.StatusCode != 200 {
		return provider.BadSessionErr
	}
	return nil
}

func (n nightbotService) URL() string {
	var j struct{
		Url string
	}
	req, _ := http.NewRequest("GET", base + "/auth/twitch", nil)
	n.client.DoJSON(req, &j)
	return j.Url
}

func (n nightbotService) ResolveSession(vals url.Values) (string, error) {
	jreq, _ := json.Marshal(map[string]interface{}{
		"code": vals.Get("code"),
		"state": vals.Get("state"),
	})

	var jresp struct{
		AccessToken string
	}

	req, err := http.NewRequest("POST", base + "/auth/twitch", bytes.NewReader(jreq))
	if err != nil {
		return "", err
	}

	_, err = n.client.DoJSON(req, &jresp)
	if err != nil {
		return "", err
	}

	return jresp.AccessToken, nil
}

func (n nightbotService) SetCommand(name, msg string) error {
	panic("implement me")
}

func (n nightbotService) DeleteCommand(name string) error {
	panic("implement me")
}


