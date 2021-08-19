package utils

import (
	"io"
	"os"
	"strings"
)

const credFileName = "cred.txt"

func GetCred(prov string) (string, bool) {
	f, err := os.OpenFile(credFileName, os.O_RDONLY, 0o667)

	if err != nil {
		return "", false
	}

	content, _ := io.ReadAll(f)
	l := strings.Split(string(content), "\n")

	if l[0] != prov {
		return "", false
	}

	return l[1], true
}

func SetCred(prov, cred string) {
	f, err := os.OpenFile(credFileName, os.O_WRONLY|os.O_CREATE, 0o667)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	f.WriteString(prov + "\n")
	f.WriteString(cred + "\n")
}
