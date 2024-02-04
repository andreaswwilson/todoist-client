package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	BaseURL = "https://api.todoist.com/rest/v2"
)

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("token missing")
	}
	return &Client{
		BaseURL: BaseURL,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}, nil
}

func init() {
	file, err := os.OpenFile("todoistClient.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetReportCaller(true)
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	// log.SetLevel(log.DebugLevel)

}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	if c.apiKey == "" {
		return fmt.Errorf("api token missing")
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	log.WithFields(log.Fields{
		"URL":     req.URL,
		"Method":  req.Method,
		"Payload": req.Body,
	}).Debug("Sending request")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	log.WithFields(log.Fields{
		"Status": res.Status,
		"Body":   string(body),
	}).Debug("Response")
	if err != nil {
		return err
	}
	// Restore the io.ReadCloser to is original state
	res.Body = io.NopCloser(bytes.NewBuffer(body))

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d, response: %s", res.StatusCode, body)
	}

	if v == nil {
		return nil
	}
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}
