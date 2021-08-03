package provider

import (
	"fmt"
	"net/url"
)

func init() {
	Register(new(Dummy))
}

type Dummy struct {}

func (d Dummy) Name() string{
	return "dummy"
}

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
	fmt.Printf("set %s to %s\n", name, msg)
	return nil
}

func (d Dummy) DeleteCommand(name string) error {
	return nil
}



