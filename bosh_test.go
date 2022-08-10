package gogobosh_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"

	. "github.com/cloudfoundry-community/gogobosh"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	mux           *http.ServeMux
	server        *httptest.Server
	fakeUAAServer *httptest.Server
)

type MockRoute struct {
	Method   string
	Endpoint string
	Output   string
	Redirect string
}

func setup(authType string) {
	setupMockRoute(MockRoute{}, authType)
}

func setupMockRoute(mock MockRoute, authType string) {
	setupMockRoutes([]MockRoute{mock}, authType)
}

func setupMockRoutes(mockEndpoints []MockRoute, authType string) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	fakeUAAServer = FakeUAAServer()
	m := martini.New()
	m.Use(render.Renderer())
	r := martini.NewRouter()
	for _, mock := range mockEndpoints {
		method := mock.Method
		endpoint := mock.Endpoint
		output := mock.Output
		redirect := mock.Redirect
		if method == "GET" {
			if redirect != "" {
				r.Get(endpoint, func(r render.Render) {
					r.Redirect(redirect)
				})
			} else {
				r.Get(endpoint, func() string {
					return output
				})
			}
		} else if method == "POST" {
			r.Post(endpoint, func() string {
				return output
			})
		} else if method == "DELETE" {
			r.Delete(endpoint, func() (int, string) {
				return 204, output
			})
		} else if method == "PUT" {
			if redirect != "" {
				r.Put(endpoint, func(r render.Render) {
					r.Redirect(redirect)
				})
			} else {
				r.Put(endpoint, func() (int, string) {
					return 200, output
				})
			}
		}
	}
	if authType != "uaa" {
		r.Get("/info", func(r render.Render) {
			r.JSON(200, &Info{
				Name:    "bosh-lite",
				UUID:    "2daf673a-9755-4b4f-aa6d-3632fbed8019",
				Version: "1.3126.0 (00000000)",
				User:    "admin",
				CPI:     "warden_cpi",
				UserAuthentication: UserAuthentication{
					Type: "basic",
				},
			})
		})
	} else {
		r.Get("/info", func(r render.Render) {
			r.JSON(200, &Info{
				Name:    "bosh-lite",
				UUID:    "2daf673a-9755-4b4f-aa6d-3632fbed8012",
				Version: "1.3126.0 (00000000)",
				User:    "admin",
				CPI:     "warden_cpi",
				UserAuthentication: UserAuthentication{
					Type: "uaa",
					Options: struct {
						URL string `json:"url"`
					}{fakeUAAServer.URL},
				},
			})
		})
	}

	m.Action(r.Handle)
	mux.Handle("/", m)
}

func FakeUAAServer() *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	m := martini.New()
	m.Use(render.Renderer())
	r := martini.NewRouter()
	count := 1
	r.Post("/oauth/token", func(r render.Render) {
		if count == 3 {
			r.JSON(401, map[string]interface{}{
				"error":             "invalid_token",
				"error_description": "Invalid refresh token (expired)",
			})
		} else {
			r.JSON(200, map[string]interface{}{
				"token_type":    "bearer",
				"access_token":  "foobar" + strconv.Itoa(count),
				"refresh_token": "barfoo",
				"expires_in":    3,
			})
		}
		count = count + 1
	})
	r.NotFound(func() string { return "" })
	m.Action(r.Handle)
	mux.Handle("/", m)
	return server
}

func teardown() {
	server.Close()
	fakeUAAServer.Close()
}
