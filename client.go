package twilgo

import (
	"net/http"
	"net/url"
	"strings"
)

// Twilio stores basic information important for connecting to the
// twilio.com REST api such as AccountSid and AuthToken.
type Twilio struct {
	AccountSid string
	AuthToken  string
	BaseURL    string
	HTTPClient *http.Client
}

// Exception is a representation of a twilio exception.
type Exception struct {
	Status   int    `json:"status"`    // HTTP specific error code
	Message  string `json:"message"`   // HTTP error message
	Code     int    `json:"code"`      // Twilio specific error code
	MoreInfo string `json:"more_info"` // Additional info from Twilio
}

// NewTwilioClient create a new Twilio struct.
func NewTwilioClient(accountSid, authToken string) *Twilio {
	return NewTwilioClientCustomHTTP(accountSid, authToken, nil)
}

// NewTwilioClientCustomHTTP create a new Twilio client, optionally using a
// custom http.Client
func NewTwilioClientCustomHTTP(accountSid, authToken string, HTTPClient *http.Client) *Twilio {
	twilioURL := "https://api.twilio.com/2010-04-01" // Should this be moved into a constant?

	if HTTPClient == nil {
		HTTPClient = http.DefaultClient
	}

	return &Twilio{accountSid, authToken, twilioURL, HTTPClient}
}

func (t *Twilio) post(formValues url.Values, twilioURL string) (*http.Response, error) {
	req, reqErr := http.NewRequest("POST", twilioURL, strings.NewReader(formValues.Encode()))
	if reqErr != nil {
		return nil, reqErr
	}
	req.SetBasicAuth(t.AccountSid, t.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := t.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	return client.Do(req)
}

func (t *Twilio) get(twilioURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", twilioURL, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(t.AccountSid, t.AuthToken)

	client := t.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	return client.Do(req)
}
