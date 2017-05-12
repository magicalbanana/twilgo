package twilgo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// CallbackParameters these are the paramters to use when you want Twilio to
// use callback urls.  See http://www.twilio.com/docs/api/rest/making-calls
// for more info.
type CallbackParameters struct {
	URL                  string // Required
	Method               string // Optional
	FallbackURL          string // Optional
	FallbackMethod       string // Optional
	StatusCallback       string // Optional
	StatusCallbackMethod string // Optional
	SendDigits           string // Optional
	IfMachine            string // False, Continue or Hangup; http://www.twilio.com/docs/errors/21207
	Timeout              int    // Optional
	Record               bool   // Optional
}

// VoiceResponse contains the details about successful voice calls.
type VoiceResponse struct {
	Sid            string   `json:"sid"`
	DateCreated    string   `json:"date_created"`
	DateUpdated    string   `json:"date_updated"`
	ParentCallSid  string   `json:"parent_call_sid"`
	AccountSid     string   `json:"account_sid"`
	To             string   `json:"to"`
	ToFormatted    string   `json:"to_formatted"`
	From           string   `json:"from"`
	FromFormatted  string   `json:"from_formatted"`
	PhoneNumberSid string   `json:"phone_number_sid"`
	Status         string   `json:"status"`
	StartTime      string   `json:"start_time"`
	EndTime        string   `json:"end_time"`
	Duration       int      `json:"duration"`
	Price          *float32 `json:"price,omitempty"`
	Direction      string   `json:"direction"`
	AnsweredBy     string   `json:"answered_by"`
	APIVersion     string   `json:"api_version"`
	Annotation     string   `json:"annotation"`
	ForwardedFrom  string   `json:"forwarded_from"`
	GroupSid       string   `json:"group_sid"`
	CallerName     string   `json:"caller_name"`
	URI            string   `json:"uri"`
}

// DateCreatedAsTime returns VoiceResponse.DateCreated as a time.Time object
// instead of a string.
func (vr *VoiceResponse) DateCreatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.DateCreated)
}

// DateUpdatedAsTime returns VoiceResponse.DateUpdated as a time.Time object
// instead of a string.
func (vr *VoiceResponse) DateUpdatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.DateUpdated)
}

// StartTimeAsTime returns VoiceResponse.StartTime as a time.Time object
// instead of a string.
func (vr *VoiceResponse) StartTimeAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.StartTime)
}

// EndTimeAsTime returns VoiceResponse.EndTime as a time.Time object instead
// of a string.
func (vr *VoiceResponse) EndTimeAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, vr.EndTime)
}

// NewCallbackParameters returns a CallbackParameters type with the specified
// url and CallbackParameters.Timeout set to 60.
func NewCallbackParameters(url string) *CallbackParameters {
	return &CallbackParameters{URL: url, Timeout: 60}
}

// CallWithURLCallbacks place a voice call with a list of callbacks specified.
func (t *Twilio) CallWithURLCallbacks(from, to string, callbackParameters *CallbackParameters) (*VoiceResponse, *Exception, error) {
	formValues := url.Values{}
	formValues.Set("From", from)
	formValues.Set("To", to)
	formValues.Set("Url", callbackParameters.URL)

	// Optional values
	if callbackParameters.Method != "" {
		formValues.Set("Method", callbackParameters.Method)
	}
	if callbackParameters.FallbackURL != "" {
		formValues.Set("FallbackURL", callbackParameters.FallbackURL)
	}
	if callbackParameters.FallbackMethod != "" {
		formValues.Set("FallbackMethod", callbackParameters.FallbackMethod)
	}
	if callbackParameters.StatusCallback != "" {
		formValues.Set("StatusCallback", callbackParameters.StatusCallback)
	}
	if callbackParameters.StatusCallbackMethod != "" {
		formValues.Set("StatusCallbackMethod", callbackParameters.StatusCallbackMethod)
	}
	if callbackParameters.SendDigits != "" {
		formValues.Set("SendDigits", callbackParameters.SendDigits)
	}
	if callbackParameters.IfMachine != "" {
		formValues.Set("IfMachine", callbackParameters.IfMachine)
	}
	if callbackParameters.Timeout != 0 {
		formValues.Set("Timeout", strconv.Itoa(callbackParameters.Timeout))
	}
	if callbackParameters.Record {
		formValues.Set("Record", "true")
	} else {
		formValues.Set("Record", "false")
	}

	return t.voicePost(formValues)
}

// CallWithApplicationCallbacks place a voice call with an ApplicationSid specified.
func (t *Twilio) CallWithApplicationCallbacks(from, to, applicationSid string) (*VoiceResponse, *Exception, error) {
	formValues := url.Values{}
	formValues.Set("From", from)
	formValues.Set("To", to)
	formValues.Set("ApplicationSid", applicationSid)

	return t.voicePost(formValues)
}

// This is a private method that has the common bits for making a voice call.
func (t *Twilio) voicePost(formValues url.Values) (*VoiceResponse, *Exception, error) {
	var voiceResponse *VoiceResponse
	var exception *Exception
	twilioURL := t.BaseURL + "/Accounts/" + t.AccountSid + "/Calls.json"

	res, err := t.post(formValues, twilioURL)
	if err != nil {
		return voiceResponse, exception, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = decoder.Decode(exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return voiceResponse, exception, err
	}

	voiceResponse = new(VoiceResponse)
	err = decoder.Decode(voiceResponse)
	return voiceResponse, exception, err
}
