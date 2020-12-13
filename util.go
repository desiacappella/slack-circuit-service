package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func epanic(err error, msg string) {
	panic(fmt.Sprintf("%v, %s", err, msg))
}

// This token needs to be acquired after authenticating with the app, on the app's home page
func getToken(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Unable to find user token")
	}
	return strings.TrimSpace(string(b))
}
