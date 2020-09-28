package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// HTTPError is an error returned by a failed API call
type HTTPError struct {
	StatusCode int
	RequestURL *url.URL
	Message    string
}

func (err HTTPError) Error() string {
	if err.Message != "" {
		return fmt.Sprintf("HTTP %d: %s (%s)", err.StatusCode, err.Message, err.RequestURL)
	}
	return fmt.Sprintf("HTTP %d (%s)", err.StatusCode, err.RequestURL)
}

type Client struct {
	http        *http.Client
	baseUrl     string
	accessToken string
}

// ClientOption represents an argument to NewClient
type ClientOption = func(http.RoundTripper) http.RoundTripper

// NewHTTPClient initializes an http.Client
func NewHTTPClient(opts ...ClientOption) *http.Client {
	tr := http.DefaultTransport
	for _, opt := range opts {
		tr = opt(tr)
	}
	return &http.Client{Transport: tr}
}

// NewClient initializes a Client
func NewClient(opts ...ClientOption) *Client {
	client := &Client{http: NewHTTPClient(opts...)}
	return client
}

// NewClientFromHTTP takes in an http.Client instance
func NewClientFromHTTP(httpClient *http.Client) *Client {
	client := &Client{http: httpClient}
	return client
}

func (c Client) FullUrl(p string) string {
	url := strings.TrimRight(c.baseUrl, " /") + p
	return url
}

func (c Client) FullWebsocketUrl(p string) string {
	url := strings.TrimRight(c.baseUrl, " /")
	parts := strings.SplitN(url, ":", 2)
	schema := strings.ToLower(parts[0])
	if schema == "http" {
		return "ws:" + parts[1] + p
	} else {
		return "wss:" + parts[1] + p
	}
}
// REST performs a REST request and parses the response.
func (c Client) REST(method string, p string, body io.Reader, data interface{}) error {
	url := strings.TrimRight(c.baseUrl, " /") + p
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return c.rest(req, data)
}

func (c Client) rest(req *http.Request, data interface{}) error {
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		return handleHTTPError(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	return nil
}

func handleHTTPError(resp *http.Response) error {
	httpError := HTTPError{
		StatusCode: resp.StatusCode,
		RequestURL: resp.Request.URL,
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		httpError.Message = err.Error()
		return httpError
	}

	var parsedBody struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &parsedBody); err == nil {
		httpError.Message = parsedBody.Message
	}

	return httpError
}
