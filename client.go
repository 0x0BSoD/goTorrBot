package transmissionRPC

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Url      string
	User     string
	Password string
	Token    string
	Http     *http.Client
}

// ===========================================
// NewClient return client instance with token
func NewClient(url, user, password string) (*Client, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := Client{
		Url:      url,
		User:     user,
		Password: password,
		Http:     &http.Client{Transport: tr},
	}

	resp, err := client.getResponse("GET", "/", nil)
	if err != nil {
		return nil, fmt.Errorf("error during getting token: %s", err)
	}

	client.Token = resp.Header.Get("X-Transmission-Session-Id")

	return &client, nil
}

// =====================================================================================================================
// PUBLIC
// =====================================================================================================================

// ===========================================
// Get - raw request, return []byte and error
func (c *Client) Get(endpoint string) ([]byte, error) {

	resp, err := c.getResponse("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error during executing request: %s", err)
	}

	// update token
	c.Token = resp.Header.Get("X-Transmission-Session-Id")

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error during read response body: %s", err)
	}

	return bodyByte, nil
}

// ===========================================
// Post - raw request, return []byte and error
func (c *Client) Post(endpoint string, body []byte) ([]byte, error) {

	resp, err := c.getResponse("POST", endpoint, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		resp, err := c.getResponse("GET", "/", nil)
		if err != nil {
			return nil, fmt.Errorf("error during getting token: %s", err)
		}
		c.Token = resp.Header.Get("X-Transmission-Session-Id")

		// try again
		return c.Post(endpoint, body)
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("error during post request: %s", resp.Status)
	}

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error during read response body: %s", err)
	}

	return bodyByte, nil
}

// =====================================================================================================================
// PRIVATE
// =====================================================================================================================

// ===========================================
// getResponse return http.Response and error [PRIVATE]
func (c *Client) getResponse(method, endpoint string, body []byte) (*http.Response, error) {
	urlStr := c.Url + endpoint

	// check full URL
	_url, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error during parsing request URL: %s", err)
	}

	// read body if present
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, _url.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error during creation of request: %s", err)
	}

	// auth
	if c.User != "" {
		req.SetBasicAuth(c.User, c.Password)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Transmission-Session-Id", c.Token)

	resp, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()

	return resp, nil
}

func (c *Client) apiCall(p *Request) ([]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return []byte{}, err
	}

	data, err := c.Post("/", b)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
