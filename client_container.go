package containerstation

import (
	"context"
	"fmt"

	"encoding/json"
	"net/http"
)

/*
	//TODO Is there any better way to do this outside of unmarshalling into interface{} or implementing unmarshal?
	//Most CS endpoints have dynamically changing JSON which is a pain in a language like Go, but this one in particular
	//has several return types instead of just two and the API docs give no indication as to all of the return types
	//that are possible. The HTTP status code for every endpoint is also useless, it seems to be 200 no matter what.
	var tmpBuf [10]byte
	if _, err = io.ReadFull(resp.Body, tmpBuf[:]); err != nil {
		return nil, err
	}

	var list []*Container
	var cse orError
	dec := json.NewDecoder(io.MultiReader(bytes.NewReader(tmpBuf[:]), resp.Body))
	if bytes.Contains(tmpBuf[:], []byte("error")) {
		err = dec.Decode(&cse)
	} else {
		err = dec.Decode(&list)
	}

	if err != nil {
		return nil, err
	}

	return list, checkEmbededCSError(cse)
*/

func (c Client) ListContainers(ctx context.Context) ([]*Container, error) {
	const apiEndpoint = `/containerstation/api/v1/container`

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
		return nil, fmt.Errorf("containerstation: Unexpected http status code from container list (%d)", resp.StatusCode)
	}

	var list []*Container
	eu := newEU(&list)

	if err = json.NewDecoder(resp.Body).Decode(eu); err != nil {
		return nil, err
	}

	return list, checkCSError(eu.orError)
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
