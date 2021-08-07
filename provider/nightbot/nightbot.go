package nightbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/memsdm05/nplink/provider"
	"github.com/memsdm05/nplink/utils"
	"io"
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
	Id        string    `json:"_id"`
}

func (n nightbotService) Name() string {
	return "nightbot"
}

func (n *nightbotService) Init() {
	n.client = utils.NewClient()
	n.client.Header.Set("referer", "https://nightbot.tv/")
	n.client.Header.Set("origin", "https://nightbot.tv")
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
		n.commands[c.Name] = &c
		fmt.Printf("%+v\n", c)
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
	req.Header.Set("Content-Type", "application/json")
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
	c := n.commands[name]
	c.Message = msg

	j, _ := json.Marshal(*c)
	fmt.Println(string(j))
	req, _ := http.NewRequest("PUT",
		base+"/1/commands/"+c.Id,
		bytes.NewReader(j))
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
	return err
}

func (n nightbotService) DeleteCommand(name string) error {
	panic("implement me")
}
