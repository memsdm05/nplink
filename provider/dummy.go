package provider

import "net/url"

type Dummy struct {}

func (d Dummy) Init() {
	panic("implement me")
}

func (d Dummy) Session(session string) error {
	return nil
}

func (d Dummy) URL() string {
	return "https://example.com"
}

func (d Dummy) ResolveSession(vals url.Values) (string, error) {
	return "foo", nil
}

func (d Dummy) SetCommand(name, msg string, extra map[string]interface{}) {}

func (d Dummy) DeleteCommand() {}