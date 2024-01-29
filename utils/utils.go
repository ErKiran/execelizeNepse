package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"
)

type Client struct {
	httpClient *http.Client
	BaseURL    *url.URL
	UserAgent  string
	Headers    string
}

func NewClient(httpClient *http.Client, apiURL string, auth string) *Client {
	baseURL, err := url.Parse(apiURL)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if err != nil {
		fmt.Println("unable to parse URL", err)
		return nil
	}

	if httpClient == nil {
		httpClient = &http.Client{Transport: tr}
	}

	c := &Client{
		httpClient: httpClient,
		UserAgent:  defaultUserAgent,
		BaseURL:    baseURL,
		Headers:    auth,
	}
	return c
}

// Do sends an API request and returns the API response. The API response is JSON
// decoded and stored in the value pointed to by v.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	if req == nil {
		return nil, errors.New("request can't be nil")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil && err != io.EOF {
		return nil, err
	}
	return resp, nil
}

// NewRequest creates an API request. The given URL is relative to the Client's
// BaseURL.
func (c *Client) NewRequest(method, urlStr string, body interface{}, opts ...interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.Reader
	if body != nil {
		// Serialize the body to JSON only if it's not nil and can be serialized.
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		fmt.Println("error in NewRequest http", err)
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func checkResponse(r *http.Response) error {
	status := r.StatusCode
	if status >= 200 && status <= 299 {
		return nil
	}

	return fmt.Errorf("request failed with status %d", status)
}
