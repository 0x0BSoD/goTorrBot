package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	Url      string
	User     string
	Password string
	Token    string
	Debug    bool
	Http     *http.Client
}

// Parameters for new client
type Parameters struct {
	Url      string
	User     string
	Password string
	Debug    bool
}

// ===========================================
// Get - raw request, return []byte and error
func (c *Client) get(endpoint string) ([]byte, error) {

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
func (c *Client) post(endpoint string, body []byte) ([]byte, error) {

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
		return c.post(endpoint, body)
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("error during post request: %s", resp.Status)
	}

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error during read response body: %s", err)
	}

	return bodyByte, nil
}

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

type Response struct {
	Result    string                 `json:"result"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Tag       int                    `json:"tag,omitempty"`
}

//
type Request struct {
	Method    string       `json:"method"`
	Arguments ReqArguments `json:"arguments,omitempty"`
	Tag       int          `json:"tag,omitempty"`
	Format    string       `json:"format,omitempty"`
}

type ReqArguments struct {
	Fields            []string `json:"fields,omitempty"`
	IDs               []int    `json:"ids,omitempty"`
	FileName          string   `json:"filename,omitempty"`
	DownloadDir       string   `json:"download-dir,omitempty"`
	MetaInfo          string   `json:"metainfo,omitempty"`
	Paused            bool     `json:"paused,omitempty"`
	PeerLimit         int      `json:"peer-limit,omitempty"`
	BandwidthPriority int      `json:"bandwidth-priority,omitempty"`
	FilesWanted       []string `json:"files-wanted,omitempty"`
	FilesUnwanted     []string `json:"files-unwanted,omitempty"`
	PriorityHigh      []int    `json:"priority-high,omitempty"`
	PriorityLow       []int    `json:"priority-low,omitempty"`
	PriorityNormal    []int    `json:"priority-normal,omitempty"`
	DeleteLocalData   bool     `json:"delete-local-data"`
	Path              string   `json:"path"`
}

func (c *Client) ApiCall(p *Request) ([]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return []byte{}, err
	}

	data, err := c.post("/", b)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
