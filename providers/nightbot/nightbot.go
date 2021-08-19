package nightbot

import (
	"bytes"
	"encoding/json"
	"github.com/memsdm05/nplink/internal/provider"
	"github.com/memsdm05/nplink/internal/utils"
	"net/http"
	"net/url"
	"time"
)

func init() {
	provider.Register(new(nightbotService))
}

const base = "https://api.nightbot.tv"

type nightbotService struct {
	client   *utils.Client
	commands map[string]*command
}

type command struct {
	CoolDown  int       `json:"coolDown"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"createdAt"`
	Message   string    `json:"message"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserLevel string    `json:"userLevel"`
	Alias     string    `json:"alias,omitempty"`
	Id        string    `json:"_id"`
}

func (n nightbotService) Name() string {
	return "nightbot"
}

func (n *nightbotService) Init() {
	n.client = utils.NewClient()
	n.client.Header.Set("referer", "https://nightbot.tv/")
	n.client.Header.Set("origin", "https://nightbot.tv")
	n.client.Header.Set("Content-Type", "application/json")
	n.commands = make(map[string]*command)
}

func (n *nightbotService) Session(session string) error {
	n.client.Header.Set("authorization", "Session "+session)

	var j struct {
		Commands []command
	}

	req, _ := http.NewRequest("GET", base+"/1/commands", nil)
	resp, err := n.client.DoJSON(req, &j)
	if resp.StatusCode != 200 || err != nil {
		return provider.BadSessionErr
	}

	for _, c := range j.Commands {
		x := c // have to copy c or else pointer complications
		n.commands[c.Name] = &x
	}

	return nil
}

func (n nightbotService) URL() string {
	var j struct {
		Url string
	}
	req, _ := http.NewRequest("GET", base+"/auth/twitch", nil)
	n.client.DoJSON(req, &j)
	return j.Url
}

func (n nightbotService) ResolveSession(vals url.Values) (string, error) {
	jreq, _ := json.Marshal(map[string]interface{}{
		"code":  vals.Get("code"),
		"state": vals.Get("state"),
	})

	var jresp struct {
		AccessToken string
	}

	req, err := http.NewRequest("POST", base+"/auth/twitch", bytes.NewReader(jreq))
	if err != nil {
		return "", err
	}

	_, err = n.client.DoJSON(req, &jresp)
	if err != nil {
		return "", err
	}

	return jresp.AccessToken, nil
}

func (n *nightbotService) makeCommand(name, msg string) error{
	c := &command{
		CoolDown: 5,
		Message: msg,
		Name: name,
		UserLevel: "everyone",
	}

	jreq, _ := json.Marshal(*c)
	var jresp struct{
		Command command
	}

	req, _ := http.NewRequest("POST",
		base+"/1/commands",
		bytes.NewReader(jreq))
	_, err := n.client.DoJSON(req, &jresp)

	if err == nil {
		n.commands[name] = &jresp.Command
	}


	return err
}

func (n *nightbotService) setCommand(c *command) error{
	j, _ := json.Marshal(*c)

	req, _ := http.NewRequest("PUT",
		base+"/1/commands/"+c.Id,
		bytes.NewReader(j))
	_, err := n.client.Do(req)

	return err
}

func (n nightbotService) SetCommand(name, msg string) error {
	if c, ok := n.commands[name]; ok {
		if c.Message == msg {
			return nil
		}
		c.Message = msg
		return n.setCommand(c)
	} else {
		return n.makeCommand(name, msg)
	}
}

func (n nightbotService) DeleteCommand(name string) error {
	panic("implement me")
}
