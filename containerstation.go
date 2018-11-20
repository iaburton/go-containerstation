package containerstation

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"net/http/cookiejar"
)

//Client is the type that implements the Container Station client v1 API found here
//http://qnap-dev.github.io/container-station-api/system.html.
//Testing has noted differences between what the API document says and the types returned by
//the various endpoints.
type Client struct {
	baseURL *string
	hc      *http.Client
}

//NewClient returns a newly initialized Client using hc as the underlying http client and baseURL for requests.
//If hc is nil than a new http client is created using similar defaults as the http.DefaultClient but with a CookieJar.
//Note that it does not use the http.DefaultClient and hc must have a CookieJar set, or one will be set.
func NewClient(baseURL string, hc *http.Client) Client {
	if hc == nil {
		//This function never returns an error
		jar, _ := cookiejar.New(nil)
		hc = &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          10,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}
	} else if hc.Jar == nil {
		//This function never returns an error
		jar, _ := cookiejar.New(nil)
		hc.Jar = jar
	}

	return Client{baseURL: &baseURL, hc: hc}
}

//Login authenticates the Client with the NAS and should normally be run before other methods.
//It takes a context for http request propagation, as well as the username and password of the account
//with access to Container Station.
func (c Client) Login(ctx context.Context, user, pass string) (*LoginResponse, error) {
	const apiEndpoint = `/containerstation/api/v1/login`

	uv := make(url.Values, 2)
	uv.Set("username", user)
	uv.Set("password", pass)
	req, err := http.NewRequest(http.MethodPost, *c.baseURL+apiEndpoint, strings.NewReader(uv.Encode()))
	i, err := c.boilerplateHTTP(ctx, req, err, loginResponseType, func(r *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	})

	//We know this will be the return type and is only nil on error
	lr, _ := i.(*LoginResponse)
	return lr, err
}

//LoginRefresh presumably refreshes the session belonging to this Client. Unfortunately the Container Station
//API docs don't explicitly state what this does, when or how often it should be called.
func (c Client) LoginRefresh(ctx context.Context) (*LoginResponse, error) {
	const apiEndpoint = `/containerstation/api/v1/login_refresh`

	req, err := http.NewRequest(http.MethodGet, *c.baseURL+apiEndpoint, nil)
	i, err := c.boilerplateHTTP(ctx, req, err, loginResponseType, nil)

	lr, _ := i.(*LoginResponse)
	return lr, err
}

//Logout invalidates the session of the Client with the NAS if it has one. It should be called when work
//with the Client is finished.
func (c Client) Logout(ctx context.Context) (*LogoutResponse, error) {
	const apiEndpoint = `/containerstation/api/v1/logout`

	req, err := http.NewRequest(http.MethodPut, *c.baseURL+apiEndpoint, nil)
	i, err := c.boilerplateHTTP(ctx, req, err, logoutResponseType, nil)

	lr, _ := i.(*LogoutResponse)
	return lr, err
}
