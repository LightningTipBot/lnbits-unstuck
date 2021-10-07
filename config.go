package main

import (
	"net/http"
	"net/url"
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
	// RetrySleepDuration required sleep duration between retries while checking if LNBits is online (default 3 seconds)
	RetrySleepDuration = 0
	// client is the http client used to connect to LNBits
	client = http.Client{}
	// lnBitsEndpoint is the parsed lnbits api endpoint
	lnBitsEndpoint *url.URL
)
