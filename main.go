package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sfreiberg/gotwilio"
)

type configurationData struct {
	From       string `json:"from"`
	To         string `json:"to"`
	AccountSID string `json:"account_sid"`
	AuthToken  string `json:"auth_token"`
}

var configFileLocations = [...]string{
	"./config.json",
	"~/.textme.json",
	"/etc/textme.json",
	"/etc/textme/config.json",
}

func main() {
	if len(os.Args) == 0 {
		fmt.Println("No message included")
		os.Exit(1)
	}
	message := os.Args[1]
	var config configurationData

	// get config
	opened := false
	for _, fp := range configFileLocations {
		abspath, _ := filepath.Abs(fp)
		jsonFile, err := os.Open(abspath)
		if err != nil {
			// try to move to next file
			continue
		}
		defer jsonFile.Close()
		data, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(data, &config)
		opened = true
		break
	}
	if !opened {
		fmt.Println("Could not find configuration file")
		os.Exit(1)
	}

	// from config
	from := config.From
	to := config.To
	accountSid := config.AccountSID
	authToken := config.AuthToken
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)
	_, _, err := twilio.SendSMS(from, to, message, "", "")
	if err != nil {
		fmt.Println("Failed to send")
		os.Exit(1)
	}
	fmt.Println("Success")
}
