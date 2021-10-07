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

func main() {
	setLogger()
	checkConfiguration()
	startWalletMonitoring()
}

// startWalletMonitoring will start the wallet monitoring
func startWalletMonitoring() {
	var err error
	lnBitsEndpoint, err = url.Parse(fmt.Sprintf("%s/api/v1/wallet", LNBitsHost))
	if err != nil {
		log.Errorf("Invalid LnBitsHost: %s", lnBitsEndpoint.String())
		panic(err)
	}
	client.Timeout = time.Second * time.Duration(TimeoutDuration)
	for {
		err := retry(3, RetrySleepDuration, checkOnlineStatus)
		if err != nil {
			log.Errorf("LNBits timout. Attempting restart.")
			log.Errorf(err.Error())
			tryTelegramNotification()
			tryRestartLnBitsCommand()
		}
		time.Sleep(time.Second * time.Duration(SleepDuration))
	}
}
func checkOnlineStatus() error {
	r := &http.Request{Header: http.Header{}, URL: lnBitsEndpoint}
	r.Header.Set("X-API-KEY", LNBitsWalletXApiKey)
	_, err := client.Do(r)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return err
	}
	log.Info("LNBits is up.")
	return nil

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
	// set default sleep after connection retry
	if RetrySleepDuration == 0 {
		RetrySleepDuration = 3
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

func retry(times, sleep int, f func() error) error {
	count := 0
	var err error
	for count < times {
		count++
		err = f()
		if err != nil {
			log.Println("LNBits seems to be down. Retrying connection")
			if count != times {
				time.Sleep(time.Duration(sleep) * time.Second)
			}
			continue
		} else {
			return nil
		}
	}
	return err
}

// setLogger will initialize the log format
func setLogger() {
	log.SetLevel(log.InfoLevel)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
}
