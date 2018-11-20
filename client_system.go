package containerstation

import (
	"context"
	"fmt"
	"net/http"
)

//SystemInformation returns information on the system that is running container station.
func (c Client) SystemInformation(ctx context.Context) (*SystemInformation, error) {
	const apiEndpoint = `/containerstation/api/v1/system`

	req, err := http.NewRequest(http.MethodGet, *c.baseURL+apiEndpoint, nil)
	i, err := c.boilerplateHTTP(ctx, req, err, sysInformationType, nil)

	si, _ := i.(*SystemInformation)
	return si, err
}

//ResourceUsage returns resource usage information on the system that is running container station.
func (c Client) ResourceUsage(ctx context.Context) (*ResourceUsage, error) {
	const apiEndpoint = `/containerstation/api/v1/system/resource`

	req, err := http.NewRequest(http.MethodGet, *c.baseURL+apiEndpoint, nil)
	i, err := c.boilerplateHTTP(ctx, req, err, resUsageType, nil)

	ru, _ := i.(*ResourceUsage)
	return ru, err
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
	i, err := c.boilerplateHTTP(ctx, req, err, npUsedType, nil)

	var toReturn bool
	if np, ok := i.(*networkPortUsed); ok {
		toReturn = np.Used
	}
	return toReturn, err
}
