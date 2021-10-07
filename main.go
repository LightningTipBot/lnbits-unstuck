package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

const (
	// TelegramBotApiKey optional telegram bot api key. If you provide Telegram bot API key, you'll get notifications
	TelegramBotApiKey = ""
	// TelegramMessage optional telegram bot that all TelegramUsers will receive
	TelegramMessage = ""
	// LNBitsHost required url to your LNBits
	LNBitsHost = "https://lnbits.com"
	// LNBitsWalletXApiKey required LNBits X-API-Key (from any wallet you want to check)
	LNBitsWalletXApiKey = ""
	// RestartLnBitsCommand optional bash command to restart your LNBits
	RestartLnBitsCommand = ""
)

var (
	// TelegramUsers optional telegram users you want to notify, if LNBits API runs into timeout
	TelegramUsers = []string{}
	// TimeoutDuration required timeout duration for the wallet API call (default 5 seconds)
	TimeoutDuration = 0
	// SleepDuration required interval you want to check you wallet using the LNBits API (default 30 seconds)
	SleepDuration = 0
	// SleepAfterRestartDuration required interval you want sleet after the LNBits restart (default 60 seconds)
	SleepAfterRestartDuration = 0
)

func main() {
	setLogger()
	checkConfiguration()
	startWalletMonitoring()
}

// startWalletMonitoring will start the wallet monitoring
func startWalletMonitoring() {
	lnBitsEndpoint, err := url.Parse(fmt.Sprintf("%s/api/v1/wallet", LNBitsHost))
	if err != nil {
		log.Errorf("Invalid LnBitsHost")
		panic(err)
	}
	c := http.Client{Timeout: time.Second * time.Duration(TimeoutDuration)}
	for {
		r := &http.Request{Header: http.Header{}, URL: lnBitsEndpoint}
		r.Header.Set("X-API-KEY", LNBitsWalletXApiKey)
		_, err := c.Do(r)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Errorf("LNBits timout. Attempting restart.")
			log.Errorf(err.Error())
			tryTelegramNotification()
			tryRestartLnBitsCommand()
		} else {
			log.Info("LNBits is up.")
		}
		time.Sleep(time.Second * time.Duration(SleepDuration))
	}
}

// tryTelegramNotification will try to notify TelegramUsers about the occurred timeout
func tryTelegramNotification() {
	if len(TelegramBotApiKey) > 0 && len(TelegramUsers) > 0 {
		for _, user := range TelegramUsers {
			r := strings.NewReader(fmt.Sprintf(`{"chat_id": "%s", "text": "%s", "disable_notification": false}`, user, TelegramMessage))
			http.Post(fmt.Sprintf("https://api.telegram.org/%s/sendMessage", TelegramBotApiKey),
				"application/json", r)
		}
	}
}

// tryRestartLnBitsCommand will try to restart LNBits using the provided command
func tryRestartLnBitsCommand() {
	if len(RestartLnBitsCommand) > 0 {
		exec.Command(RestartLnBitsCommand)
	}
}

// checkConfiguration will check if the provided configuration is valid
func checkConfiguration() {
	if len(LNBitsHost) == 0 {
		panic(fmt.Errorf("please provide a LNBits host (LNBitsHost)"))
	}
	if len(LNBitsWalletXApiKey) == 0 {
		panic(fmt.Errorf("please provide a LNBits API key (LNBitsWalletXApiKey)"))
	}
	// set default timeout
	if TimeoutDuration == 0 {
		TimeoutDuration = 5
	}
	// set default sleep
	if SleepDuration == 0 {
		SleepDuration = 30
	}
	// set default sleep after restart
	if SleepAfterRestartDuration == 0 {
		SleepAfterRestartDuration = 60
	}
	if len(TelegramBotApiKey) > 0 && len(TelegramUsers) == 0 {
		panic(fmt.Errorf("please provide the telegram users id's you want to notify (TelegramUsers)"))
	}
	if len(TelegramBotApiKey) > 0 && len(TelegramUsers) > 0 && TelegramMessage == "" {
		panic(fmt.Errorf("please provide a valid TelegramMessage"))
	}
	if (len(TelegramBotApiKey) == 0 || len(TelegramUsers) == 0) && len(RestartLnBitsCommand) == 0 {
		panic(fmt.Errorf("please provide either a telegram API key (LnBitsHost) and user id's (TelegramUsers) or a valid RestartLnBitsCommand. Otherwise this application is useless"))
	}

}

// setLogger will initialize the log format
func setLogger() {
	log.SetLevel(log.InfoLevel)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
}
