package twilgo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// SmsResponse is returned after a text/sms message is posted to Twilio
type SmsResponse struct {
	Sid         string   `json:"sid"`
	DateCreated string   `json:"date_created"`
	DateUpdate  string   `json:"date_updated"`
	DateSent    string   `json:"date_sent"`
	AccountSid  string   `json:"account_sid"`
	To          string   `json:"to"`
	From        string   `json:"from"`
	MediaURL    string   `json:"media_url"`
	Body        string   `json:"body"`
	Status      string   `json:"status"`
	Direction   string   `json:"direction"`
	APIVersion  string   `json:"api_version"`
	Price       *float32 `json:"price,omitempty"`
	URL         string   `json:"uri"`
}

// DateCreatedAsTime returns SmsResponse.DateCreated as a time.Time object
// instead of a string.
func (s *SmsResponse) DateCreatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, s.DateCreated)
}

// DateUpdateAsTime returns SmsResponse.DateUpdate as a time.Time object
// instead of a string.
func (s *SmsResponse) DateUpdateAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, s.DateUpdate)
}

// DateSentAsTime returns SmsResponse.DateSent as a time.Time object instead
// of a string.
func (s *SmsResponse) DateSentAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, s.DateSent)
}

// SendSMS uses Twilio to send a text message. See
// http://www.twilio.com/docs/api/rest/sending-sms for more information.
func (t *Twilio) SendSMS(from, to, body, statusCallback, applicationSid string) (*SmsResponse, *Exception, error) {
	formValues := initFormValues(to, body, "", statusCallback, applicationSid)
	formValues.Set("From", from)

	return t.sendMessage(formValues)
}

// SendSMSWithCopilot uses Twilio Copilot to send a text message.
// See https://www.twilio.com/docs/api/rest/sending-messages-copilot
func (t *Twilio) SendSMSWithCopilot(messagingServiceSid, to, body, statusCallback, applicationSid string) (*SmsResponse, *Exception, error) {
	formValues := initFormValues(to, body, "", statusCallback, applicationSid)
	formValues.Set("MessagingServiceSid", messagingServiceSid)

	return t.sendMessage(formValues)
}

// SendMMS uses Twilio to send a multimedia message.
func (t *Twilio) SendMMS(from, to, body, mediaURL, statusCallback, applicationSid string) (*SmsResponse, *Exception, error) {
	formValues := initFormValues(to, body, mediaURL, statusCallback, applicationSid)
	formValues.Set("From", from)

	return t.sendMessage(formValues)
}

// Core method to send message
func (t *Twilio) sendMessage(formValues url.Values) (*SmsResponse, *Exception, error) {
	twilioURL := t.BaseURL + "/Accounts/" + t.AccountSid + "/Messages.json"

	res, postErr := t.post(formValues, twilioURL)
	if postErr != nil {
		return nil, nil, postErr
	}
	defer res.Body.Close()

	responseBody, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, nil, readErr
	}

	// if the status code does not return created check for the exception that
	// was returned.
	if res.StatusCode != http.StatusCreated {
		exception := &Exception{}
		unMarshalErr := json.Unmarshal(responseBody, exception)

		if unMarshalErr != nil {
			return nil, nil, unMarshalErr
		}

		return nil, exception, nil
	}

	smsResponse := &SmsResponse{}
	unMarshalErr := json.Unmarshal(responseBody, smsResponse)
	if unMarshalErr != nil {
		return nil, nil, unMarshalErr
	}

	return smsResponse, nil, nil
}

// Form values initialization
func initFormValues(to, body, mediaURL, statusCallback, applicationSid string) url.Values {
	formValues := url.Values{}

	formValues.Set("To", to)
	formValues.Set("Body", body)

	if mediaURL != "" {
		formValues.Set("MediaURL", mediaURL)
	}

	if statusCallback != "" {
		formValues.Set("StatusCallback", statusCallback)
	}

	if applicationSid != "" {
		formValues.Set("ApplicationSid", applicationSid)
	}

	return formValues
}
