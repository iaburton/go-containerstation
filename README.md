

# containerstation
`import "github.com/iaburton/containerstation"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
Package containerstation implements a Go API client for QNAP's ContainerStation API found here.

<a href="http://qnap-dev.github.io/container-station-api/index.html">http://qnap-dev.github.io/container-station-api/index.html</a>

The package name is just 'containerstation' not go-containerstation as the url/repo implies, and importing
under a shorter alias such a cstation is recommended.
Please note this package is a work in progress; more endpoints, tests and comments need to be added.
A licence will be added and the package opensourced as it gets closer to an initial version/release.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [type Client](#Client)
  * [func NewClient(baseURL string, hc *http.Client) Client](#NewClient)
  * [func (c Client) DownloadTLSCertificate(ctx context.Context, file string, perm os.FileMode) (err error)](#Client.DownloadTLSCertificate)
  * [func (c Client) ExportTLSCertificate(ctx context.Context, w io.Writer) error](#Client.ExportTLSCertificate)
  * [func (c Client) GetContainer(ctx context.Context, ctype cType, id string) (*Container, error)](#Client.GetContainer)
  * [func (c Client) ListContainers(ctx context.Context) ([]*Container, error)](#Client.ListContainers)
  * [func (c Client) Login(ctx context.Context, user, pass string) (*LoginResponse, error)](#Client.Login)
  * [func (c Client) LoginRefresh(ctx context.Context) (*LoginResponse, error)](#Client.LoginRefresh)
  * [func (c Client) Logout(ctx context.Context) (*LogoutResponse, error)](#Client.Logout)
  * [func (c Client) NetworkPort(ctx context.Context, proto protocol, port int) (bool, error)](#Client.NetworkPort)
  * [func (c Client) RemoveContainer(ctx context.Context, ctype cType, id string) (*Container, error)](#Client.RemoveContainer)
  * [func (c Client) ResourceUsage(ctx context.Context) (*ResourceUsage, error)](#Client.ResourceUsage)
  * [func (c Client) RestartContainer(ctx context.Context, ctype cType, id string) (*Container, error)](#Client.RestartContainer)
  * [func (c Client) StartContainer(ctx context.Context, ctype cType, id string) (*Container, error)](#Client.StartContainer)
  * [func (c Client) StopContainer(ctx context.Context, ctype cType, id string) (*Container, error)](#Client.StopContainer)
  * [func (c Client) SystemInformation(ctx context.Context) (*SystemInformation, error)](#Client.SystemInformation)
* [type Container](#Container)
* [type LoginResponse](#LoginResponse)
* [type LogoutResponse](#LogoutResponse)
* [type ResourceUsage](#ResourceUsage)
* [type SystemInformation](#SystemInformation)

#### <a name="pkg-examples">Examples</a>
* [NewClient](#example_NewClient)

#### <a name="pkg-files">Package files</a>
[client_container.go](/src/github.com/iaburton/containerstation/client_container.go) [client_system.go](/src/github.com/iaburton/containerstation/client_system.go) [client_tls.go](/src/github.com/iaburton/containerstation/client_tls.go) [containerstation.go](/src/github.com/iaburton/containerstation/containerstation.go) [json_reflect.go](/src/github.com/iaburton/containerstation/json_reflect.go) [json_types.go](/src/github.com/iaburton/containerstation/json_types.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    DOCKER cType = iota + 1
    LXC
)
```
``` go
const (
    //TCP enum to be used with client.NetworkPort
    TCP protocol = iota
    //UDP enum to be used with client.NetworkPort
    UDP
)
```




## <a name="Client">type</a> [Client](/src/target/containerstation.go?s=912:973#L28)
``` go
type Client struct {
    // contains filtered or unexported fields
}
```
Client is the type that implements the Container Station client v1 API found here
<a href="http://qnap-dev.github.io/container-station-api/system.html">http://qnap-dev.github.io/container-station-api/system.html</a>.
Testing has noted differences between what the API document says and the types returned by
the various endpoints.







### <a name="NewClient">func</a> [NewClient](/src/target/containerstation.go?s=1312:1366#L36)
``` go
func NewClient(baseURL string, hc *http.Client) Client
```
NewClient returns a newly initialized Client using hc as the underlying http client and baseURL for requests.
If hc is nil than a new http client is created using similar defaults as the http.DefaultClient but with a CookieJar.
Note that it does not use the http.DefaultClient and hc must have a CookieJar set, or one will be set.





### <a name="Client.DownloadTLSCertificate">func</a> (Client) [DownloadTLSCertificate](/src/target/client_tls.go?s=73:175#L10)
``` go
func (c Client) DownloadTLSCertificate(ctx context.Context, file string, perm os.FileMode) (err error)
```



### <a name="Client.ExportTLSCertificate">func</a> (Client) [ExportTLSCertificate](/src/target/client_tls.go?s=470:546#L27)
``` go
func (c Client) ExportTLSCertificate(ctx context.Context, w io.Writer) error
```



### <a name="Client.GetContainer">func</a> (Client) [GetContainer](/src/target/client_container.go?s=1623:1716#L80)
``` go
func (c Client) GetContainer(ctx context.Context, ctype cType, id string) (*Container, error)
```



### <a name="Client.ListContainers">func</a> (Client) [ListContainers](/src/target/client_container.go?s=69:142#L10)
``` go
func (c Client) ListContainers(ctx context.Context) ([]*Container, error)
```



### <a name="Client.Login">func</a> (Client) [Login](/src/target/containerstation.go?s=2273:2358#L67)
``` go
func (c Client) Login(ctx context.Context, user, pass string) (*LoginResponse, error)
```
Login authenticates the Client with the NAS and should normally be run before other methods.
It takes a context for http request propagation, as well as the username and password of the account
with access to Container Station.




### <a name="Client.LoginRefresh">func</a> (Client) [LoginRefresh](/src/target/containerstation.go?s=3071:3144#L85)
``` go
func (c Client) LoginRefresh(ctx context.Context) (*LoginResponse, error)
```
LoginRefresh presumably refreshes the session belonging to this Client. Unfortunately the Container Station
API docs don't explicitly state what this does, when or how often it should be called.




### <a name="Client.Logout">func</a> (Client) [Logout](/src/target/containerstation.go?s=3537:3605#L97)
``` go
func (c Client) Logout(ctx context.Context) (*LogoutResponse, error)
```
Logout invalidates the session of the Client with the NAS if it has one. It should be called when work
with the Client is finished.




### <a name="Client.NetworkPort">func</a> (Client) [NetworkPort](/src/target/client_system.go?s=1558:1646#L63)
``` go
func (c Client) NetworkPort(ctx context.Context, proto protocol, port int) (bool, error)
```
NetworkPort reports whether or not a given protocol and port are being used.
Protocol must be one of TCP or UDP constants defined in this package.
If an invalid protocol or port is set an error is returned.




### <a name="Client.RemoveContainer">func</a> (Client) [RemoveContainer](/src/target/client_container.go?s=2363:2459#L97)
``` go
func (c Client) RemoveContainer(ctx context.Context, ctype cType, id string) (*Container, error)
```
TODO mention nil error and "zero" value Container struct on success




### <a name="Client.ResourceUsage">func</a> (Client) [ResourceUsage](/src/target/client_system.go?s=594:668#L21)
``` go
func (c Client) ResourceUsage(ctx context.Context) (*ResourceUsage, error)
```
ResourceUsage returns resource usage information on the system that is running container station.




### <a name="Client.RestartContainer">func</a> (Client) [RestartContainer](/src/target/client_container.go?s=2120:2217#L92)
``` go
func (c Client) RestartContainer(ctx context.Context, ctype cType, id string) (*Container, error)
```



### <a name="Client.StartContainer">func</a> (Client) [StartContainer](/src/target/client_container.go?s=1784:1879#L84)
``` go
func (c Client) StartContainer(ctx context.Context, ctype cType, id string) (*Container, error)
```



### <a name="Client.StopContainer">func</a> (Client) [StopContainer](/src/target/client_container.go?s=1953:2047#L88)
``` go
func (c Client) StopContainer(ctx context.Context, ctype cType, id string) (*Container, error)
```



### <a name="Client.SystemInformation">func</a> (Client) [SystemInformation](/src/target/client_system.go?s=157:239#L10)
``` go
func (c Client) SystemInformation(ctx context.Context) (*SystemInformation, error)
```
SystemInformation returns information on the system that is running container station.




## <a name="Container">type</a> [Container](/src/target/json_types.go?s=1883:2324#L63)
``` go
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
    Type      cType    `json:"type"`
}
```
Container is the JSON returned when getting basic information on a container
or list of containers.










## <a name="LoginResponse">type</a> [LoginResponse](/src/target/json_types.go?s=138:311#L10)
``` go
type LoginResponse struct {
    Anonymous bool   `json:"anonymous"`
    IsAdmin   bool   `json:"isAdmin"`
    Time      string `json:"loginime"`
    Username  string `json:"username"`
}
```
LoginResponse is the JSON returned for a login or login refresh.










## <a name="LogoutResponse">type</a> [LogoutResponse](/src/target/json_types.go?s=373:438#L18)
``` go
type LogoutResponse struct {
    Username string `json:"username"`
}
```
LogoutResponse is the JSON returned for a logout request.










## <a name="ResourceUsage">type</a> [ResourceUsage](/src/target/json_types.go?s=1397:1777#L48)
``` go
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
```
ResourceUsage is the JSON returned for a resource usage request.










## <a name="SystemInformation">type</a> [SystemInformation](/src/target/json_types.go?s=515:1328#L23)
``` go
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
```
SystemInformation is the JSON returned for a system information request.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
