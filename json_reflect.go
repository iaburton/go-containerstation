package containerstation

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"net/http"
)

var (
	//used in conjunction with the method below to cut down on serious boilerplate/very similar redundant code
	loginResponseType  = reflect.TypeOf(LoginResponse{})
	logoutResponseType = reflect.TypeOf(LogoutResponse{})
	sysInformationType = reflect.TypeOf(SystemInformation{})
	resUsageType       = reflect.TypeOf(ResourceUsage{})
	containerType      = reflect.TypeOf(Container{})
	npUsedType         = reflect.TypeOf(networkPortUsed{})
	containerSliceType = reflect.TypeOf([]*Container{})
)

func (c Client) boilerplateHTTP(ctx context.Context, req *http.Request, err error, tp reflect.Type, onReq func(*http.Request)) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if onReq != nil {
		onReq(req)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//TODO any point in checking beyond 200?
	//Seems containerstation only returns 200 regardless of what actually happened.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("containerstation: Unexpected http status code (%d) from (%s) at (%s)", resp.StatusCode, req.URL.Host, req.URL.Path)
	}

	instance := reflect.New(tp).Interface()
	eu := newEU(instance)
	if err := json.NewDecoder(resp.Body).Decode(eu); err != nil {
		return nil, err
	}

	return instance, checkCSError(eu.orError)
}
