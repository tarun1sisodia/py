package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// SendSMS sends an SMS with the provided message to the given phone number using Twilio's Messaging API.
func SendSMS(phone string, message string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	from := os.Getenv("TWILIO_FROM") // ensure you set this in your .env file
	if accountSid == "" || authToken == "" || from == "" {
		return fmt.Errorf("Twilio credentials (TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_FROM) not set")
	}

	// Twilio Messaging API endpoint
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	// Prepare data for the POST request
	msgData := url.Values{}
	msgData.Set("To", phone)
	msgData.Set("From", from)
	msgData.Set("Body", message)
	payload := strings.NewReader(msgData.Encode())

	req, err := http.NewRequest("POST", urlStr, payload)
	if err != nil {
		return err
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Twilio API error: %s", string(bodyBytes))
	}

	log.Printf("SMS sent to %s successfully.", phone)
	return nil
}
