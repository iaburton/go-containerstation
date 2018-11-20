package containerstation

import (
	"context"

	"net/http"
)

func (c Client) ListContainers(ctx context.Context) ([]*Container, error) {
	const apiEndpoint = `/containerstation/api/v1/container`

	req, err := http.NewRequest(http.MethodGet, *c.baseURL+apiEndpoint, nil)
	i, err := c.boilerplateHTTP(ctx, req, err, containerSliceType, nil)

	list, _ := i.(*[]*Container)
	return *list, err
}

type cType int

const (
	DOCKER cType = iota
	LXC
)

func (c cType) String() string {
	switch c {
	case DOCKER:
		return "docker"
	case LXC:
		return "lxc"
	default:
		return "<unknown container type>"
	}
}

func (c cType) isValid() bool {
	switch c {
	case DOCKER, LXC:
		return true
	default:
		return false
	}
}

// func (c Client) GetContainer(ctx context.Context) (*Container, error) {

// }
