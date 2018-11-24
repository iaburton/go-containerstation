package containerstation

import (
	"context"
	"fmt"

	"net/http"
)

func (c Client) ListContainers(ctx context.Context) ([]*Container, error) {
	const apiEndpoint = `/containerstation/api/v1/container`

	req, err := http.NewRequest(http.MethodGet, *c.baseURL+apiEndpoint, nil)
	i, err := c.boilerplateHTTP(ctx, req, err, containerSliceType, nil)

	var toReturn []*Container
	if list, ok := i.(*[]*Container); ok {
		toReturn = *list
	}
	return toReturn, err
}

type cType int

const (
	DOCKER cType = iota + 1
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

func (c *cType) UnmarshalJSON(p []byte) error {
	if l := len(p); l == 5 || l == 8 {
		switch string(p) {
		case `"docker"`:
			*c = DOCKER
		case `"lxc"`:
			*c = LXC
		default:
			return fmt.Errorf("cType: Unknown container type (%s)", p)
		}
		return nil
	}

	return fmt.Errorf("cType: Unknown container type")
}

func (c Client) containerAction(ctx context.Context, ctype cType, id, finalEndpoint, method string) (*Container, error) {
	const apiEndpoint = `/containerstation/api/v1/container/%s/%s`

	if !ctype.isValid() {
		return nil, fmt.Errorf("containerstation: Bad container type (%d) in call to Container", ctype)
	}

	req, err := http.NewRequest(method, *c.baseURL+fmt.Sprintf(apiEndpoint, ctype, id)+finalEndpoint, nil)
	i, err := c.boilerplateHTTP(ctx, req, err, containerType, nil)

	ct, _ := i.(*Container)
	return ct, err
}

func (c Client) GetContainer(ctx context.Context, ctype cType, id string) (*Container, error) {
	return c.containerAction(ctx, ctype, id, "", http.MethodGet)
}

func (c Client) StartContainer(ctx context.Context, ctype cType, id string) (*Container, error) {
	return c.containerAction(ctx, ctype, id, "/start", http.MethodPut)
}

func (c Client) StopContainer(ctx context.Context, ctype cType, id string) (*Container, error) {
	return c.containerAction(ctx, ctype, id, "/stop", http.MethodPut)
}

func (c Client) RestartContainer(ctx context.Context, ctype cType, id string) (*Container, error) {
	return c.containerAction(ctx, ctype, id, "/restart", http.MethodPut)
}

//TODO mention nil error and "zero" value Container struct on success
func (c Client) RemoveContainer(ctx context.Context, ctype cType, id string) (*Container, error) {
	return c.containerAction(ctx, ctype, id, "", http.MethodDelete)
}
