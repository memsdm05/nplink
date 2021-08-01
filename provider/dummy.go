package provider

import "net/url"

func init() {
	Register("dummy", new(Dummy))
}

type Dummy struct {}

func (d Dummy) Init() {}

func (d Dummy) Session(session string) error {
	return nil
}

func (d Dummy) URL() string {
	return "https://example.com"
}

func (d Dummy) ResolveSession(vals url.Values) (string, error) {
	return "foobar", nil
}

func (d Dummy) SetCommand(name, msg string) error {
	return nil
}

func (d Dummy) DeleteCommand(name string) error {
	return nil
}



