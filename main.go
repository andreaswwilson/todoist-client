package todoist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	BaseURL = "https://api.todoist.com/rest/v2"
)

type Client struct {
	BaseURL    string
	ApiKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("token missing")
	}
	return &Client{
		BaseURL: BaseURL,
		ApiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}, nil
}

func init() {
	log.SetReportCaller(true)
	// log.SetLevel(log.DebugLevel)

}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	if c.ApiKey == "" {
		return fmt.Errorf("api token missing")
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))

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

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
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
