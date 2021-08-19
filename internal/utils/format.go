package utils

import (
	"fmt"
	"strings"
)

const (
	StartDelimiter = '{'
	EndDelimiter   = '}'
)

type FMap map[string]string

func (f FMap) Set(key, value string) {
	f[key] = value
}

func (f FMap) Setf(key, format string, values ...interface{}) {
	f.Set(key, fmt.Sprintf(format, values...))
}

func (f FMap) SetFunc(key string, strfunc func() string) {
	f.Set(key, strfunc())
}

type token struct {
	r       bool
	content string
}

type FormatString struct {
	tokens []token
}

func NewFormatString(inp string) *FormatString {
	ret := new(FormatString)
	ret.tokens = buildTokens(inp)

	return ret
}

func buildTokens(inp string) []token {
	ret := make([]token, 0)
	o := 0
	f := false
	nxt := rune(0)

	for i, char := range inp {
		if i < len(inp)-1 {
			nxt = rune(inp[i+1])
		}

		if char == StartDelimiter && char != nxt {
			if i-o > 0 {
				ret = append(ret, token{
					r:       false,
					content: inp[o:i],
				})
			}
			o = i + 1
			f = true
		}

		if char == EndDelimiter && f {
			ret = append(ret, token{
				r:       true,
				content: inp[o:i],
			})
			o = i + 1
			f = false
		}

		if i == len(inp)-1 {
			if i-o > 0 {
				ret = append(ret, token{
					r:       false,
					content: inp[o : i+1],
				})
			}
		}
	}

	return ret
}

func (f *FormatString) Format(fmap FMap) string {
	var sb strings.Builder

	for _, token := range f.tokens {
		if token.r {

			if v, ok := fmap[token.content]; ok {
				sb.WriteString(v)
			} else {
				sb.WriteString(string(StartDelimiter) + token.content + string(EndDelimiter))
			}

		} else {
			sb.WriteString(token.content)
		}
	}

	return sb.String()
}

func (f *FormatString) UnmarshalText(text []byte) error {
	f.tokens = buildTokens(string(text))
	return nil
}
