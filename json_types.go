package containerstation

import (
	"bytes"
	"encoding/json"
	"fmt"
)

//LoginResponse is the JSON returned for a login or login refresh.
type LoginResponse struct {
	Anonymous bool   `json:"anonymous"`
	IsAdmin   bool   `json:"isAdmin"`
	Time      string `json:"loginime"`
	Username  string `json:"username"`
}

//LogoutResponse is the JSON returned for a logout request.
type LogoutResponse struct {
	Username string `json:"username"`
}

//SystemInformation is the JSON returned for a system information request.
type SystemInformation struct {
	CPUCore   int      `json:"cpuCore"`
	CPUThread int      `json:"cpuThread"`
	Features  []string `json:"features"`
	Gpu       struct {
		CsMode          bool          `json:"cs_mode"`
		Device          []interface{} `json:"device"`
		DriverInstalled bool          `json:"driver_installed"`
	} `json:"gpu"`
	GpuDriver   bool   `json:"gpuDriver"`
	Hostname    string `json:"hostname"`
	Machine     string `json:"machine"`
	NeedRestart bool   `json:"needRestart"`
	Processor   string `json:"processor"`
	Status      string `json:"status"`
	Version     struct {
		DockerVersion string `json:"dockerVersion"`
		Firmware      string `json:"firmware"`
		LxcVersion    string `json:"lxcVersion"`
		Qpkg          string `json:"qpkg"`
		Web           string `json:"web"`
	} `json:"version"`
}

//ResourceUsage is the JSON returned for a resource usage request.
type ResourceUsage struct {
	CPU    string `json:"cpu_usage"`
	Memory struct {
		Buffers        int `json:"buffers"`
		Cached         int `json:"cached"`
		Percent        int `json:"percent"`
		PercentBuffers int `json:"percent_buffers"`
		PercentCached  int `json:"percent_cached"`
		Total          int `json:"total"`
		Used           int `json:"used"`
	} `json:"memory_usage"`
}

//Container is the JSON returned when getting basic information on a container
//or list of containers.
type Container struct {
	CPU       float64  `json:"cpu"`
	ID        string   `json:"id"`
	Image     string   `json:"image"`
	ImageID   string   `json:"imageID"`
	Ipaddress []string `json:"ipaddress"`
	Memory    int      `json:"memory"`
	Name      string   `json:"name"`
	Rx        int      `json:"rx"`
	State     string   `json:"state"`
	TCPPort   []int    `json:"tcpPort"`
	Tx        int      `json:"tx"`
	Type      string   `json:"type"`
}

type networkPortUsed struct {
	Used bool `json:"used"`
}

type errorUnmarshaler struct {
	expected interface{}
	orError
}

func newEU(expected interface{}) *errorUnmarshaler {
	return &errorUnmarshaler{expected: expected}
}

func (eu *errorUnmarshaler) UnmarshalJSON(p []byte) error {
	if len(p) > 10 && bytes.Equal(p[4:11], []byte(`"error"`)) {
		return json.Unmarshal(p, &eu.orError)
	}

	return json.Unmarshal(p, eu.expected)
}

type orError struct {
	Error *csError `json:"error,omitempty"`
}

type csError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *csError) Error() string {
	return fmt.Sprintf("ContainerStation error (%d): %s", e.Code, e.Message)
}

func checkCSError(oe orError) error {
	if oe.Error == nil {
		return nil
	}
	return oe.Error
}
