package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/net/websocket"
)

type Frame struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	SecondsLeft float64 `json:"seconds_left"`
}

const origin = "http://localhost/"

var wsRe = regexp.MustCompile(`wss://wss.redditmedia.com/thebutton\?h=[0-9a-f]*&e=[0-9a-f]*`)

func getWsUrl() (string, error) {
	resp, err := http.Get("https://reddit.com/r/thebutton")
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	match := wsRe.Find(body)
	if match == nil {
		return "", errors.New("Could not find websocket URL")
	}

	return string(match), nil
}

func main() {
	url, err := getWsUrl()
	if err != nil {
		log.Fatal(err)
	}

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	for {
		var data Frame
		websocket.JSON.Receive(ws, &data)

		fmt.Println(data.Payload.SecondsLeft)
	}
}
