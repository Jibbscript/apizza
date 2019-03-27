package dawg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	// WarnigStatus is the status code dominos serves use for a warning
	WarnigStatus = 1
	// FailureStatus  is the status code dominos serves use for a failure
	FailureStatus = -1
	// OkStatus  is the status code dominos serves use to signify no problems
	OkStatus = 0

	// DefaultLang is the package language variable
	DefaultLang = "en"

	host = "order.dominos.com"
)

var (
	// Warnings is a package switch for turning warings on or off
	Warnings = false

	cli = &http.Client{
		Timeout: time.Duration(10 * time.Second),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	errCodes = map[int]string{
		FailureStatus: "Failure -1",
		WarnigStatus:  "Warning 1",
		OkStatus:      "Ok 0",
	}
)

func dominosErr(resp []byte) error {
	e := &DominosError{}
	if err := e.init(resp); err != nil {
		return err
	}

	if e.IsOk() {
		return nil
	}
	return e
}

// DominosError represents an error sent back by the dominos servers
type DominosError struct {
	Status      int
	StatusItems []statusItem
	Order       struct {
		Status      int
		StatusItems []statusItem
	}
	Msg     string
	fullErr map[string]interface{}
}

type statusItem struct {
	Code      string
	Message   string
	PulseCode int
	PulseText string
}

// Init initializes the error from json data.
func (err *DominosError) init(jsonData []byte) error {
	err.fullErr = map[string]interface{}{}

	if e := json.Unmarshal(jsonData, &err.fullErr); e != nil {
		return e
	}
	return mapstructure.Decode(err.fullErr, err)
}

func (err *DominosError) Error() string {
	var errmsg string

	for _, item := range err.StatusItems {
		errmsg += fmt.Sprintf("Dominos %s:\n", item.Code)
	}
	for _, item := range err.Order.StatusItems {
		if item.Code != "" {
			errmsg += fmt.Sprintf("    Code: '%s'", item.Code)
		}
		if item.Message != "" {
			errmsg += fmt.Sprintf(":\n        %s\n", item.Message)
		} else if item.PulseText != "" {
			errmsg += fmt.Sprintf("    PulseCode %d:\n        %s", item.PulseCode, item.PulseText)
		} else {
			errmsg += "\n"
		}
	}
	return errmsg
}

// IsWarning returns true when the error sent by dominos is a warning
func (err *DominosError) IsWarning() bool {
	return err.Status == WarnigStatus
}

// IsFailure returns true if the error that dominos sent back prevents the
// system from working
func (err *DominosError) IsFailure() bool {
	return err.Status == FailureStatus
}

// IsOk returns true is the error is not a failure else returns false
func (err *DominosError) IsOk() bool {
	return err.Status != FailureStatus
}

// PrintData prints the raw data. (just for testing... you shouldn't be seeing this)
func (err *DominosError) PrintData() {
	for key, value := range err.fullErr {
		if key == "Order" {
			fmt.Println("Order:")
			for k, v := range value.(map[string]interface{}) {
				fmt.Printf("  %s: %v\n", k, v)
			}
		} else {
			fmt.Printf("%s: %v\n", key, value)
		}
	}
}

func get(path string, params URLParam) ([]byte, error) {
	if params == nil {
		params = &Params{}
	}
	return send(&http.Request{
		Method: "GET",
		Host:   host,
		Proto:  "HTTP/1.1",
		Header: make(http.Header),
		URL: &url.URL{
			Scheme:   "https",
			Host:     host,
			Path:     path,
			RawQuery: params.Encode(),
		},
	})
}

func post(path string, data []byte) ([]byte, error) {
	return send(&http.Request{
		Method: "POST",
		Host:   host,
		Proto:  "HTTP/1.1",
		Body:   ioutil.NopCloser(bytes.NewReader(data)),
		Header: make(http.Header),
		URL: &url.URL{
			Scheme: "https",
			Host:   host,
			Path:   path,
		},
	})
}

func send(req *http.Request) ([]byte, error) {
	var buf bytes.Buffer

	req.Header.Add("User-Agent", "Dominos API Wrapper for GO - "+time.Now().String())

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return buf.Bytes(), fmt.Errorf("bad response code: %d", resp.StatusCode)
	}
	return buf.Bytes(), nil
}
