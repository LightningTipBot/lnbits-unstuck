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
	// TelegramBotApiKey optional telegram bot api key
	TelegramBotApiKey = ""
	// LnBitsHost required url to your lnbits
	LnBitsHost = "https://lnbits.com"
	// XApiKey required lnbits X-API-Key (from any wallet you want to check)
	XApiKey = ""
	// RestartLnBitsCommand optional bash command to restart your lnbits
	RestartLnBitsCommand = ""
)

var (
	// TelegramUsers optional telegram users you want to notify, if lnbits API runs into timeout
	TelegramUsers = []string{}
	// TimeoutDuration required timeout duration for the wallet API call (default 5 seconds)
	TimeoutDuration = 0
	// SleepDuration required interval you want to check you wallet using the lnbits API (default 30 seconds)
	SleepDuration = 0
)

func main() {
	checkConfiguration()
	startWalletMonitoring()
}

// startWalletMonitoring will start the wallet monitoring
func startWalletMonitoring() {
	lnBitsEndpoint, err := url.Parse(fmt.Sprintf("%s/api/v1/wallet", LnBitsHost))
	if err != nil {
		panic(err)
	}
	c := http.Client{Timeout: time.Second * time.Duration(TimeoutDuration)}
	for {
		r := &http.Request{Header: http.Header{}, URL: lnBitsEndpoint}
		r.Header.Set("X-API-KEY", XApiKey)
		_, err := c.Do(r)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			tryTelegramNotification()
			tryRestartLnBitsCommand()
		}
		time.Sleep(time.Second * time.Duration(SleepDuration))
	}
}

// tryTelegramNotification will try to notify TelegramUsers about the occurred timeout
func tryTelegramNotification() {
	if len(TelegramBotApiKey) > 0 && len(TelegramUsers) > 0 {
		for _, user := range TelegramUsers {
			r := strings.NewReader(fmt.Sprintf(`{"chat_id": "%s", "text": "LNBITS API request timeout! Trying to restart", "disable_notification": false}`, user))
			http.Post(fmt.Sprintf("https://api.telegram.org/%s/sendMessage", TelegramBotApiKey),
				"application/json", r)
		}
	}
}

// tryRestartLnBitsCommand will try to restart lnbits using the provided command
func tryRestartLnBitsCommand() {
	if len(RestartLnBitsCommand) > 0 {
		exec.Command(RestartLnBitsCommand)
	}
}

// checkConfiguration will check if the provided configuration is valid
func checkConfiguration() {
	if len(LnBitsHost) == 0 {
		panic(fmt.Errorf("please provide a lnbits host"))
	}
	if len(XApiKey) == 0 {
		panic(fmt.Errorf("please provide a lnbits API key"))
	}
	// set default timeout
	if TimeoutDuration == 0 {
		TimeoutDuration = 5
	}
	// set default sleep
	if SleepDuration == 0 {
		SleepDuration = 30
	}
	if len(TelegramBotApiKey) > 0 && len(TelegramUsers) == 0 {
		panic(fmt.Errorf("please provide the telegram users id's you want to notify"))
	}
	if (len(TelegramBotApiKey) == 0 || len(TelegramUsers) == 0) && len(RestartLnBitsCommand) == 0 {
		panic(fmt.Errorf("please provide either a telegram API key and user id's or a valid RestartLnBitsCommand. Otherwise this application is useless"))
	}
}
