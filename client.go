package gogobosh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//Client used to communicate with BOSH
type Client struct {
	config Config
}

//Config is used to configure the creation of a client
type Config struct {
	BOSHAddress       string
	Port              string
	Username          string
	Password          string
	HttpClient        *http.Client
	SkipSslValidation bool
}

// request is used to help build up a request
type request struct {
	method string
	url    string
	header map[string]string
	params url.Values
	body   io.Reader
	obj    interface{}
}

//DefaultConfig configuration for client
func DefaultConfig() *Config {
	return &Config{
		BOSHAddress:       "https://192.168.50.4:25555",
		Username:          "admin",
		Password:          "admin",
		HttpClient:        http.DefaultClient,
		SkipSslValidation: true,
	}
}

// NewClient returns a new client
func NewClient(config *Config) *Client {
	// bootstrap the config
	defConfig := DefaultConfig()

	if len(config.BOSHAddress) == 0 {
		config.BOSHAddress = defConfig.BOSHAddress
	}

	if len(config.Username) == 0 {
		config.Username = defConfig.Username
	}

	if len(config.Password) == 0 {
		config.Password = defConfig.Password
	}

	config.HttpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.SkipSslValidation,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			req.URL.Host = strings.TrimPrefix(config.BOSHAddress, req.URL.Scheme+"://")
			req.SetBasicAuth(config.Username, config.Password)
			return nil
		}}
	client := &Client{
		config: *config,
	}
	return client
}

// NewRequest is used to create a new request
func (c *Client) NewRequest(method, path string) *request {
	r := &request{
		method: method,
		url:    c.config.BOSHAddress + path,
		params: make(map[string][]string),
		header: make(map[string]string),
	}
	return r
}

// DoRequest runs a request with our client
func (c *Client) DoRequest(r *request) (*http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return nil, err
	}
	for key, value := range r.header {
		req.Header.Add(key, value)
	}
	req.SetBasicAuth(c.config.Username, c.config.Password)
	req.Header.Add("User-Agent", "gogo-bosh")
	resp, err := c.config.HttpClient.Do(req)
	return resp, err
}

// UUID return uuid
func (c *Client) UUID() string {
	info, _ := c.GetInfo()
	return info.UUID
}

// GetInfo returns BOSH Info
func (c *Client) GetInfo() (info Info, err error) {
	r := c.NewRequest("GET", "/info")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting info %v", err)
		return
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading info request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &info)
	if err != nil {
		log.Printf("Error unmarshaling info %v", err)
		return
	}
	return
}

// toHTTP converts the request to an HTTP request
func (r *request) toHTTP() (*http.Request, error) {

	// Check if we should encode the body
	if r.body == nil && r.obj != nil {
		if b, err := encodeBody(r.obj); err != nil {
			return nil, err
		} else {
			r.body = b
		}
	}

	// Create the HTTP request
	return http.NewRequest(r.method, r.url, r.body)
}

// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

// encodeBody is used to encode a request body
func encodeBody(obj interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}
