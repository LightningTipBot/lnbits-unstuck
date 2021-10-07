package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

const (
	TelegramBotApiKey    = ""
	LnBitsHost           = "https://lnbits.com"
	RestartLnBitsCommand = ""
)

var TelegramUsers = []string{}

func main() {
	if len(LnBitsHost) == 0 {
		panic(fmt.Errorf("please provide a lnbits host"))
	}
	lnBitsEndpoint, err := url.Parse(fmt.Sprintf("%s/api/v1/wallet", LnBitsHost))
	if err != nil {
		panic(err)
	}
	c := http.Client{Timeout: time.Second * 5}
	for {
		r := &http.Request{Header: http.Header{}, URL: lnBitsEndpoint}
		_, err := c.Do(r)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			if len(TelegramBotApiKey) > 0 && len(TelegramUsers) > 0 {
				for _, user := range TelegramUsers {
					r := strings.NewReader(fmt.Sprintf(`{"chat_id": "%s", "text": "LNBITS API request timeout! Trying to restart", "disable_notification": false}`, user))
					http.Post(fmt.Sprintf("https://api.telegram.org/%s/sendMessage", TelegramBotApiKey),
						"application/json", r)
				}
			}
			if len(RestartLnBitsCommand) > 0 {
				exec.Command(RestartLnBitsCommand)
			}
		}
		time.Sleep(time.Second * 30)
	}
}
