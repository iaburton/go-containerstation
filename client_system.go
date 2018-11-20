package containerstation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

//SystemInformation returns information on the system that is running container station.
func (c Client) SystemInformation(ctx context.Context) (*SystemInformation, error) {
	const apiEndpoint = `/containerstation/api/v1/system`

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
		return nil, fmt.Errorf("containerstation: Unexpected http status code from system (%d)", resp.StatusCode)
	}

	si := new(SystemInformation)
	if err := json.NewDecoder(resp.Body).Decode(si); err != nil {
		return nil, err
	}

	return si, checkCSError(si.orError)
}

//ResourceUsage returns resource usage information on the system that is running container station.
func (c Client) ResourceUsage(ctx context.Context) (*ResourceUsage, error) {
	const apiEndpoint = `/containerstation/api/v1/system/resource`

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
		return nil, fmt.Errorf("containerstation: Unexpected http status code from resource usage (%d)", resp.StatusCode)
	}

	ru := new(ResourceUsage)
	if err := json.NewDecoder(resp.Body).Decode(ru); err != nil {
		return nil, err
	}

	return ru, checkCSError(ru.orError)
}

type protocol int

const (
	//TCP enum to be used with client.NetworkPort
	TCP protocol = iota
	//UDP enum to be used with client.NetworkPort
	UDP
)

func (p protocol) String() string {
	switch p {
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	default:
		return fmt.Sprintf("<unknown protocol %d>", p)
	}
}

func (p protocol) isValid() bool {
	switch p {
	case TCP, UDP:
		return true
	default:
		return false
	}
}

//NetworkPort reports whether or not a given protocol and port are being used.
//Protocol must be one of TCP or UDP constants defined in this package.
//If an invalid protocol or port is set an error is returned.
func (c Client) NetworkPort(ctx context.Context, proto protocol, port int) (bool, error) {
	const apiEndpoint = `/containerstation/api/v1/system/port/%s/%d`

	if port < 1 || !proto.isValid() {
		return false, fmt.Errorf("containerstation: Bad port or protocol in call to NetworkPort")
	}

	req, err := http.NewRequest(http.MethodGet, *c.baseURL+fmt.Sprintf(apiEndpoint, proto, port), nil)
	if err != nil {
		return false, err
	}

	req = req.WithContext(ctx)
	resp, err := c.hc.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("containerstation: Unexpected http status code from network port (%d)", resp.StatusCode)
	}

	is := new(struct {
		Used bool `json:"used"`
		orError
	})
	if err := json.NewDecoder(resp.Body).Decode(is); err != nil {
		return false, err
	}

	return is.Used, checkCSError(is.orError)
}
