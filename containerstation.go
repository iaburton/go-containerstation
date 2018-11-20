package containerstation

import (
	"context"
	"encoding/json"
	"fmt"
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
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//TODO any point in checking beyond 200?
	//Seems containerstation only returns 200 regardless of what actually happened.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("containerstation: Unexpected http status code from login (%d)", resp.StatusCode)
	}

	lr := new(LoginResponse)
	eu := newEU(lr)
	if err := json.NewDecoder(resp.Body).Decode(eu); err != nil {
		return nil, err
	}

	return lr, checkCSError(eu.orError)
}

//LoginRefresh presumably refreshes the session belonging to this Client. Unfortunately the Container Station
//API docs don't explicitly state what this does, when or how often it should be called.
func (c Client) LoginRefresh(ctx context.Context) (*LoginResponse, error) {
	const apiEndpoint = `/containerstation/api/v1/login_refresh`

	req, err := http.NewRequest(http.MethodGet, *c.baseURL+apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("containerstation: Unexpected http status code from login_refresh (%d)", resp.StatusCode)
	}

	lr := new(LoginResponse)
	if err := json.NewDecoder(resp.Body).Decode(lr); err != nil {
		return nil, err
	}

	return lr, checkCSError(lr.orError)
}

//Logout invalidates the session of the Client with the NAS if it has one. It should be called when work
//with the Client is finished.
func (c Client) Logout(ctx context.Context) (*LogoutResponse, error) {
	const apiEndpoint = `/containerstation/api/v1/logout`

	req, err := http.NewRequest(http.MethodPut, *c.baseURL+apiEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("containerstation: Unexpected http status code from logout (%d)", resp.StatusCode)
	}

	lr := new(LogoutResponse)
	if err := json.NewDecoder(resp.Body).Decode(lr); err != nil {
		return nil, err
	}

	return lr, checkCSError(lr.orError)
}
