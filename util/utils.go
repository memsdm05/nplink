package util

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func Must(err error) {
	if err != nil {
		FatalError(err)
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
	fmt.Println("An error has occurred:\n")
	color.Red("%v", err)
	fmt.Println("\n\n< Press Enter to Exit >a")
	//fmt.Printf("%s\n\n< press enter to close >", err)
	fmt.Scanln()
	os.Exit(727) // I am very funny
}

func GetCred(prov string) (string, time.Time, error) {
	content, err := os.ReadFile("cred.txt")

	if err != nil {
		return "", time.Time{}, err
	}

	l := strings.Split(string(content), "\n")
	t, err := time.Parse(time.RFC3339, l[2])

	if err != nil {
		return "", time.Time{}, err
	}

	if l[0] != prov {
		return "", time.Time{}, errors.New("incorrect provider")
	}

	return l[1], t, nil
}

func SetCred(prov, cred string) error {
	f, err := os.OpenFile("cred.txt", os.O_WRONLY|os.O_CREATE, 0o667)
	defer f.Close()

	if err != nil {
		return err
	}

	f.WriteString(prov + "\n")
	f.WriteString(cred + "\n")
	f.WriteString(time.Now().Format(time.RFC3339) + "\n")

	return nil
}