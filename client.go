package main

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
		return nil, fmt.Errorf("error during post request: %s", err)
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

func (c *Client) apiCall(p *Parameters) ([]byte, error) {
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

// =====================================================================================================================
// RPC
// =====================================================================================================================

func (c *Client) SessionStats() (Statistics, error) {
	p := &Parameters{
		Method:    "session-stats",
		Arguments: Arguments{},
	}

	data, err := c.apiCall(p)
	if err != nil {
		return Statistics{}, err
	}

	var result Statistics
	err = json.Unmarshal(data, &result)
	if err != nil {
		return Statistics{}, err
	}

	return result, nil
}

func (c *Client) TorrentGet() {
	p := &Parameters{
		Method: "torrent-get",
		Arguments: Arguments{
			Fields: []string{"id", "name", "status", "comment",
				"error", "errorString", "isFinished",
				"leftUntilDone", "percentDone",
				"sizeWhenDone", "startDate",
				"uploadRatio", "totalSize"},
		},
	}

	data, err := c.apiCall(p)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

//func (c *Client) TorrentGetId(IDs []int) {
//	p := &Parameters{
//		Method: "torrent-get",
//		Arguments: Arguments{
//			Fields: []string{"id", "name", "status", "comment",
//				"error", "errorString", "isFinished",
//				"leftUntilDone", "percentDone", "errorString",
//				"sizeWhenDone", "startDate", "addedDate",
//				"uploadRatio", "totalSize", "peers",
//				"rateDownload", "rateUpload", "uploadRatio"},
//			IDs: IDs,
//		},
//	}
//
//	data, err := c.apiCall(p)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(string(data))
//}
