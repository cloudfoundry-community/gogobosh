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
	Username          string
	Password          string
	UAAAuth           bool
	HttpClient        *http.Client
	SkipSslValidation bool
	TokenSource       oauth2.TokenSource
	Endpoint          *Endpoint
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
		config.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			req.URL.Host = strings.TrimPrefix(config.BOSHAddress, req.URL.Scheme+"://")
			req.SetBasicAuth(config.Username, config.Password)
			req.Header.Add("User-Agent", "gogo-bosh")
			req.Header.Del("Referer")
			return nil
		}
	} else {
		ctx := getContext(*config)

		endpoint, err := getUAAEndpoint(config.BOSHAddress, oauth2.NewClient(ctx, nil))

		if err != nil {
			return nil, fmt.Errorf("Could not get api /info: %v", err)
		}

		config.Endpoint = endpoint

		authConfig, token, err := getToken(ctx, *config)

		if err != nil {
			return nil, fmt.Errorf("Error getting token: %v", err)
		}

		config.TokenSource = authConfig.TokenSource(ctx, token)
		config.HttpClient = oauth2.NewClient(ctx, config.TokenSource)

		config.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			req.URL.Host = strings.TrimPrefix(config.BOSHAddress, req.URL.Scheme+"://")
			req.Header.Add("User-Agent", "gogo-bosh")
			req.Header.Del("Referer")
			return nil
		}
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
	if err != nil {
		if strings.Contains(err.Error(), "oauth2: cannot fetch token") {
			c.refreshClient()
			resp, err = c.config.HttpClient.Do(req)
		}
	}
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

func (c *Client) refreshClient() error {
	ctx := getContext(c.config)

	authConfig, token, err := getToken(ctx, c.config)
	if err != nil {
		return fmt.Errorf("Error getting token: %v", err)
	}

	c.config.TokenSource = authConfig.TokenSource(ctx, token)
	c.config.HttpClient = oauth2.NewClient(ctx, c.config.TokenSource)
	c.config.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > 10 {
			return fmt.Errorf("stopped after 10 redirects")
		}
		req.URL.Host = strings.TrimPrefix(c.config.BOSHAddress, req.URL.Scheme+"://")
		req.Header.Add("User-Agent", "gogo-bosh")
		req.Header.Del("Referer")
		return nil
	}

	return nil
}

func getToken(ctx context.Context, config Config) (*oauth2.Config, *oauth2.Token, error) {
	authConfig := &oauth2.Config{
		ClientID: "bosh_cli",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Endpoint.URL + "/oauth/authorize",
			TokenURL: config.Endpoint.URL + "/oauth/token",
		},
	}
	token, err := authConfig.PasswordCredentialsToken(ctx, config.Username, config.Password)
	return authConfig, token, err
}

func getContext(config Config) context.Context {
	ctx := oauth2.NoContext
	if config.SkipSslValidation == false {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, config.HttpClient)
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: tr})
	}
	return ctx
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

// GetToken - returns the current token bearer
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
