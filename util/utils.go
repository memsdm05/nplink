package util

import (
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"net/url"
	"os"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

var StopRedirect = func(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

func TransposeValues(values url.Values, keys ...string) (ret url.Values) {
	for _, key := range keys { ret.Set(key, values.Get(key)) }
	return
}

func TransposeHeader(header http.Header, keys ...string) (ret http.Header) {
	for _, key := range keys { ret.Set(key, header.Get(key)) }
	return
}

func FatalError(err error)  {
	fmt.Println("An error has occurred:")
	color.Red("\n%v\n\n", err)
	fmt.Println("< Press Enter to Exit >")
	//fmt.Printf("%s\n\n< press enter to close >", err)
	fmt.Scanln()
	os.Exit(1)
}