package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	server         string
	Client         *http.Client
	defaultHeaders map[string]string
}

func NewHTTPClient(server string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		server: server,
		Client: &http.Client{
			Timeout: timeout * time.Second,
		},
		defaultHeaders: make(map[string]string),
	}
}

type RequestOption func(*http.Request)

func WithHeader(key, value string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

func WithHeaders(headers map[string]string) RequestOption {
	return func(req *http.Request) {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
}

func (c *HTTPClient) Get(path string, options ...RequestOption) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", c.server+path, nil)
	if err != nil {
		return nil, err
	}

	c.addDefaultHeaders(req)

	for _, option := range options {
		option(req)
	}

	return c.Client.Do(req)
}

func (c *HTTPClient) Post(path string, body interface{}) (resp *http.Response, err error) {
	var reader io.Reader

	switch v := body.(type) {
	case nil:
		reader = nil
	case io.Reader:
		reader = v
	case string:
		reader = strings.NewReader(v)
	case []byte:
		reader = bytes.NewReader(v)
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("JSON Marshal Fail: %v", err)
		}
		reader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest("POST", c.server+path, reader)
	if err != nil {
		return nil, err
	}

	c.addDefaultHeaders(req)

	return c.Client.Do(req)
}

func (c *HTTPClient) AddDefaultHeader(key, value string) {
	c.defaultHeaders[key] = value
}

func (c *HTTPClient) AddDefaultHeaders(headers map[string]string) {
	for key, value := range headers {
		c.defaultHeaders[key] = value
	}
}

func (c *HTTPClient) GetBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("handleHTTPResponse: Read Fail: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("handleHTTPResponse: StatusCode Error(%d) %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *HTTPClient) addDefaultHeaders(req *http.Request) {
	for key, value := range c.defaultHeaders {
		req.Header.Set(key, value)
	}
}
