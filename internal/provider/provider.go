package provider

import (
	"errors"
	"net/url"
)

var providers = make(map[string]Provider)

var BadSessionErr = errors.New("provider: bad session token")
var UnknownProviderErr = errors.New("provider: Could not find provider")

type Provider interface {
	Name() string

	// Init init provider after it's picked to reduce unused resources
	Init()

	// Session return error if session is invalid
	// Set the session and init session related attrs
	Session(session string) error

	URL() string
	ResolveSession(vals url.Values) (string, error)

	SetCommand(name, msg string) error
	DeleteCommand(name string) error
}

func Register(provider Provider) {
	providers[provider.Name()] = provider
}

func Select(name string) (Provider, error) {
	ret, ok := providers[name]

	if !ok {
		return nil, UnknownProviderErr
	}

	// clear up some memory
	for k := range providers {
		if k != name {
			delete(providers, k)
		}
	}

	ret.Init()

	return ret, nil
}
