package gogobosh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	boshhttp "github.com/cloudfoundry/bosh-utils/httpclient"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Client used to communicate with BOSH
type Client struct {
	config   Config
	Endpoint Endpoint
}

// Config is used to configure the creation of a client
type Config struct {
	BOSHAddress       string
	Username          string
	Password          string
	ClientID          string
	ClientSecret      string
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

// DefaultConfig configuration for client
func DefaultConfig() *Config {
	return &Config{
		BOSHAddress:       "https://192.168.50.4:25555", // bosh-lite default IP:PORT
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

	// Save the configured HTTP Client timeout for later
	var timeout time.Duration
	if config.HttpClient != nil {
		timeout = config.HttpClient.Timeout
	}

	// Skip TLS cert validation and respect BOSH_ALL_PROXY env var
	config.HttpClient = boshhttp.CreateDefaultClientInsecureSkipVerify()
	endpoint := &Endpoint{}

	authType, err := getAuthType(config.BOSHAddress, config.HttpClient)
	if err != nil {
		return nil, fmt.Errorf("could not get client auth type: %w", err)
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
			return nil, fmt.Errorf("could not get api /info: %w", err)
		}

		config.Endpoint = endpoint

		if config.ClientID == "" { //No ClientID? Do UAA User auth
			authConfig, token, err := getToken(ctx, *config)

			if err != nil {
				return nil, fmt.Errorf("error getting token: %w", err)
			}

			config.TokenSource = authConfig.TokenSource(ctx, token)
			config.HttpClient = oauth2.NewClient(ctx, config.TokenSource)
		} else { //Got a ClientID? Do UAA Client Auth (two-legged auth)
			authConfig := &clientcredentials.Config{
				ClientID:     config.ClientID,
				ClientSecret: config.ClientSecret,
				TokenURL:     endpoint.URL + "/oauth/token",
			}
			config.TokenSource = authConfig.TokenSource(ctx)
			config.HttpClient = authConfig.Client(ctx)
		}

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

	//Restore the timeout from the provided HTTP Client
	config.HttpClient.Timeout = timeout

	client := &Client{
		config:   *config,
		Endpoint: *endpoint,
	}

	return client, nil
}

func getAuthType(api string, httpClient *http.Client) (string, error) {
	info, err := getInfo(api, httpClient)
	return info.UserAuthentication.Type, err
}

func getInfo(api string, httpClient *http.Client) (*Info, error) {
	if api == "" {
		return &Info{}, nil
	}

	resp, err := httpClient.Get(api + "/info")
	if err != nil {
		return &Info{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	var info Info
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return &Info{}, fmt.Errorf("error unmarshalling info response: %w", err)
	}
	return &info, err
}

func getUAAEndpoint(api string, httpClient *http.Client) (*Endpoint, error) {
	if api == "" {
		return DefaultEndpoint(), nil
	}
	info, err := getInfo(api, httpClient)
	URL := info.UserAuthentication.Options.URL
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

func (c *Client) DoRequestAndUnmarshal(r *request, objPtr interface{}) error {
	resp, err := c.DoRequest(r)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	err = json.NewDecoder(resp.Body).Decode(objPtr)
	if err != nil {
		return fmt.Errorf("error unmarshalling http response: %w", err)
	}
	return nil
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
			err = c.refreshClient()
			if err != nil {
				return nil, fmt.Errorf("error refreshing UAA client: %w", err)
			}
			resp, err = c.config.HttpClient.Do(req)
		} else {
			// errors are only returned for very bad things, not 400s etc
			return nil, fmt.Errorf("error making bosh client http request: %w", err)
		}
	} else if resp.StatusCode >= 400 {
		if strings.Contains(resp.Status, "Unauthorized") {
			err = c.refreshClient()
			if err != nil {
				return nil, fmt.Errorf("error refreshing UAA client from 400: %w", err)
			}
			resp, err = c.config.HttpClient.Do(req)
		} else {
			return nil, fmt.Errorf("http %s request to %s failed with %s", req.Method, req.URL, resp.Status)
		}
	}
	return resp, err
}

// GetUUID returns the BOSH UUID
func (c *Client) GetUUID() (string, error) {
	info, err := c.GetInfo()
	if err != nil {
		return "", fmt.Errorf("error getting the UUID: %w", err)
	}
	return info.UUID, nil
}

// UUID returns the BOSH uuid
// Deprecated: Use GetUUID and check for errors
func (c *Client) UUID() string {
	uuid, _ := c.GetUUID()
	return uuid
}

// GetInfo returns BOSH Info
func (c *Client) GetInfo() (Info, error) {
	info, err := getInfo(c.config.BOSHAddress, c.config.HttpClient)
	if err != nil {
		return Info{}, err
	}
	return *info, nil
}

func (c *Client) refreshClient() error {
	// Create a new http client to avoid authentication failure when getting a new
	// token as the oauth2 client passes along the expired/revoked refresh token.
	c.config.HttpClient = &http.Client{
		Timeout: c.config.HttpClient.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.config.SkipSslValidation,
			},
		},
	}

	ctx := getContext(c.config)

	authConfig, token, err := getToken(ctx, c.config)
	if err != nil {
		return fmt.Errorf("error getting token to refresh client: %w", err)
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
	return context.WithValue(context.Background(), oauth2.HTTPClient, config.HttpClient)
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
		return "", fmt.Errorf("error getting bearer token: %w", err)
	}
	return "bearer " + token.AccessToken, nil
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
