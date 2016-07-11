package gogobosh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//Client used to communicate with BOSH
type Client struct {
	config   Config
	Endpoint Endpoint
}

//Config is used to configure the creation of a client
type Config struct {
	BOSHAddress       string
	Port              string
	Username          string
	Password          string
	UAAAuth           bool
	HttpClient        *http.Client
	SkipSslValidation bool
	TokenSource       oauth2.TokenSource
}

type Endpoint struct {
	URL string `json:"doppler_logging_endpoint"`
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

func DefaultEndpoint() *Endpoint {
	return &Endpoint{
		URL: "https://192.168.50.4:8443",
	}
}

// NewClient returns a new client
func NewClient(config *Config) (*Client, error) {
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

	endpoint := &Endpoint{}
	config.HttpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.SkipSslValidation,
			},
		},
	}
	authType, err := getAuthType(config.BOSHAddress, config.HttpClient)
	if err != nil {
		return nil, fmt.Errorf("Could not get auth type: %v", err)
	}
	if authType != "uaa" {
		config.HttpClient = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) > 10 {
					return fmt.Errorf("stopped after 10 redirects")
				}
				req.URL.Host = strings.TrimPrefix(config.BOSHAddress, req.URL.Scheme+"://")
				req.SetBasicAuth(config.Username, config.Password)
				req.Header.Add("User-Agent", "gogo-bosh")
				req.Header.Del("Referer")
				return nil
			},
		}
	} else {
		ctx := oauth2.NoContext
		if config.SkipSslValidation == false {
			ctx = context.WithValue(ctx, oauth2.HTTPClient, defConfig.HttpClient)
		} else {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: tr})
		}

		endpoint, err := getUAAEndpoint(config.BOSHAddress, oauth2.NewClient(ctx, nil))

		if err != nil {
			return nil, fmt.Errorf("Could not get api /info: %v", err)
		}

		authConfig := &oauth2.Config{
			ClientID: "cf",
			Scopes:   []string{""},
			Endpoint: oauth2.Endpoint{
				AuthURL:  endpoint.URL + "/oauth/auth",
				TokenURL: endpoint.URL + "/oauth/token",
			},
		}

		token, err := authConfig.PasswordCredentialsToken(ctx, config.Username, config.Password)
		if err != nil {
			return nil, fmt.Errorf("Error getting token: %v", err)
		}

		config.TokenSource = authConfig.TokenSource(ctx, token)
		config.HttpClient = oauth2.NewClient(ctx, config.TokenSource)
	}
	client := &Client{
		config:   *config,
		Endpoint: *endpoint,
	}

	return client, nil
}

func getAuthType(api string, httpClient *http.Client) (string, error) {
	info, err := getInfo(api, httpClient)
	return info.UserAuthenication.Type, err
}

func getInfo(api string, httpClient *http.Client) (*Info, error) {
	var (
		info Info
	)

	if api == "" {
		return &Info{}, nil
	}

	resp, err := httpClient.Get(api + "/info")
	if err != nil {
		log.Printf("Error requesting info %v", err)
		return &Info{}, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading info request %v", resBody)
		return &Info{}, err
	}
	err = json.Unmarshal(resBody, &info)
	return &info, err
}

func getUAAEndpoint(api string, httpClient *http.Client) (*Endpoint, error) {
	if api == "" {
		return DefaultEndpoint(), nil
	}
	info, err := getInfo(api, httpClient)
	URL := info.UserAuthenication.Options.URL
	return &Endpoint{URL: URL}, err
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
	defer resp.Body.Close()

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

func (c *Client) GetToken() (string, error) {
	token, err := c.config.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("Error getting bearer token: %v", err)
	}
	return "bearer " + token.AccessToken, nil
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
