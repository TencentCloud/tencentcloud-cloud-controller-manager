package clb

const (
	LoadBalanceListenerProtocolHTTP  = 1
	LoadBalanceListenerProtocolTCP   = 2
	LoadBalanceListenerProtocolUDP   = 3
	LoadBalanceListenerProtocolHTTPS = 4
)

type CreateListenerOpts struct {
	LoadBalancerPort int32   `qcloud_arg:"loadBalancerPort,required"`
	InstancePort     int32   `qcloud_arg:"instancePort,required"`
	Protocol         int     `qcloud_arg:"protocol,required"`
	ListenerName     *string `qcloud_arg:"listenerName"`
	SessionExpire    *int    `qcloud_arg:"sessionExpire"`
	HealthSwitch     *int    `qcloud_arg:"healthSwitch"`
	TimeOut          *int    `qcloud_arg:"timeOut"`
	IntervalTime     *int    `qcloud_arg:"intervalTime"`
	HealthNum        *int    `qcloud_arg:"healthNum"`
	UnhealthNum      *int    `qcloud_arg:"unhealthNum"`
	HttpHash         *int    `qcloud_arg:"httpHash"`
	HttpCode         *int    `qcloud_arg:"httpCode"`
	HttpCheckPath    *string `qcloud_arg:"httpCheckPath"`
	SSLMode          *string `qcloud_arg:"SSLMode"`
	CertId           *string `qcloud_arg:"certId"`
	CertCaId         *string `qcloud_arg:"certCaId"`
	CertCaContent    *string `qcloud_arg:"certCaContent"`
	CertCaName       *string `qcloud_arg:"certCaName"`
	CertContent      *string `qcloud_arg:"certContent"`
	CertKey          *string `qcloud_arg:"certKey"`
	CertName         *string `qcloud_arg:"certName"`
}

type CreateLoadBalancerListenersArgs struct {
	LoadBalancerId string               `qcloud_arg:"loadBalancerId,required"`
	Listeners      []CreateListenerOpts `qcloud_arg:"listeners"`
}

type CreateLoadBalancerListenersResponse struct {
	Response
	RequestId   int      `json:"requestId"`
	ListenerIds []string `json:"listenerIds"`
}

func (response CreateLoadBalancerListenersResponse) Id() int {
	return response.RequestId
}

func (client *Client) CreateLoadBalancerListeners(args *CreateLoadBalancerListenersArgs) (
	*CreateLoadBalancerListenersResponse,
	error,
) {
	response := &CreateLoadBalancerListenersResponse{}
	err := client.Invoke("CreateLoadBalancerListeners", args, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type DescribeLoadBalancerListenersArgs struct {
	LoadBalancerId   string    `qcloud_arg:"loadBalancerId,required"`
	ListenerIds      *[]string `qcloud_arg:"listenerIds"`
	Protocol         *int      `qcloud_arg:"protocol"`
	LoadBalancerPort *int32    `qcloud_arg:"loadBalancerPort"`
	Status           *int      `qcloud_arg:"status"`
}

type DescribeLoadBalancerListenersResponse struct {
	Response
	TotalCount  int        `json:"totalCount"`
	ListenerSet []Listener `json:"listenerSet"`
}

type Listener struct {
	UnListenerId     string `json:"unListenerId"`
	LoadBalancerPort int32  `json:"loadBalancerPort"`
	InstancePort     int32  `json:"instancePort"`
	Protocol         int    `json:"protocol"`
	SessionExpire    int    `json:"sessionExpire"`
	HealthSwitch     int    `json:"healthSwitch"`
	TimeOut          int    `json:"timeOut"`
	IntervalTime     int    `json:"intervalTime"`
	HealthNum        int    `json:"healthNum"`
	UnhealthNum      int    `json:"unhealthNum"`
	HttpHash         string `json:"httpHash"`
	HttpCode         int    `json:"httpCode"`
	HttpCheckPath    string `json:"httpCheckPath"`
	SSLMode          string `json:"SSLMode"`
	CertId           string `json:"certId"`
	CertCaId         string `json:"certCaId"`
	Status           int    `json:"status"`
}

func (client *Client) DescribeLoadBalancerListeners(args *DescribeLoadBalancerListenersArgs) (
	*DescribeLoadBalancerListenersResponse,
	error,
) {
	response := &DescribeLoadBalancerListenersResponse{}
	err := client.Invoke("DescribeLoadBalancerListeners", args, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type DeleteLoadBalancerListenersArgs struct {
	LoadBalancerId string   `qcloud_arg:"loadBalancerId,required"`
	ListenerIds    []string `qcloud_arg:"listenerIds,required"`
}

type DeleteLoadBalancerListenersResponse struct {
	Response
	RequestId int `json:"requestId"`
}

func (response DeleteLoadBalancerListenersResponse) Id() int {
	return response.RequestId
}

func (client *Client) DeleteLoadBalancerListeners(LoadBalancerId string, ListenerIds []string) (
	*DeleteLoadBalancerListenersResponse,
	error,
) {

	response := &DeleteLoadBalancerListenersResponse{}
	err := client.Invoke("DeleteLoadBalancerListeners", &DeleteLoadBalancerListenersArgs{
		LoadBalancerId: LoadBalancerId,
		ListenerIds:    ListenerIds,
	}, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type ModifyLoadBalancerListenerArgs struct {
	LoadBalancerId string  `qcloud_arg:"loadBalancerId,required"`
	ListenerId     string  `qcloud_arg:"listenerId,required"`
	ListenerName   *string `qcloud_arg:"listenerName"`
	SessionExpire  *int    `qcloud_arg:"sessionExpire"`
	HealthSwitch   *int    `qcloud_arg:"healthSwitch"`
	TimeOut        *int    `qcloud_arg:"timeOut"`
	IntervalTime   *int    `qcloud_arg:"intervalTime"`
	HealthNum      *int    `qcloud_arg:"healthNum"`
	UnhealthNum    *int    `qcloud_arg:"unhealthNum"`
	HttpHash       *int    `qcloud_arg:"httpHash"`
	HttpCode       *int    `qcloud_arg:"httpCode"`
	HttpCheckPath  *string `qcloud_arg:"httpCheckPath"`
	SSLMode        *string `qcloud_arg:"SSLMode"`
	CertId         *string `qcloud_arg:"certId"`
	CertCaId       *string `qcloud_arg:"certCaId"`
	CertCaContent  *string `qcloud_arg:"certCaContent"`
	CertCaName     *string `qcloud_arg:"certCaName"`
	CertContent    *string `qcloud_arg:"certContent"`
	CertKey        *string `qcloud_arg:"certKey"`
	CertName       *string `qcloud_arg:"certName"`
}

type ModifyLoadBalancerListenerResponse struct {
	Response
	RequestId int `json:"requestId"`
}

func (response ModifyLoadBalancerListenerResponse) Id() int {
	return response.RequestId
}

func (client *Client) ModifyLoadBalancerListener(args *ModifyLoadBalancerListenerArgs) (
	*ModifyLoadBalancerListenerResponse,
	error,
) {
	response := &ModifyLoadBalancerListenerResponse{}
	err := client.Invoke("ModifyLoadBalancerListener", args, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type CreateForwardLBFourthLayerListenersArgs struct {
	LoadBalancerId string                          `qcloud_arg:"loadBalancerId"`
	Listeners      []CreateFourthLayerListenerOpts `qcloud_arg:"listeners"`
}

type CreateFourthLayerListenerOpts struct {
	LoadBalancerPort int     `qcloud_arg:"loadBalancerPort"`
	Protocol         int     `qcloud_arg:"protocol"`
	ListenerName     *string `qcloud_arg:"listenerName"`
	SessionExpire    *int    `qcloud_arg:"sessionExpire"`
	HealthSwitch     *int    `qcloud_arg:"healthSwitch"`
	TimeOut          *int    `qcloud_arg:"timeOut"`
	IntervalTime     *int    `qcloud_arg:"intervalTime"`
	HealthNum        *int    `qcloud_arg:"healthNum"`
	UnhealthNum      *int    `qcloud_arg:"unhealthNum"`
	Scheduler        *string `qcloud_arg:"scheduler"`
}

type CreateForwardLBFourthLayerListenersResponse struct {
	Response
	RequestId   int      `json:"requestId"`
	ListenerIds []string `json:"listenerIds"`
}

func (response CreateForwardLBFourthLayerListenersResponse) Id() int {
	return response.RequestId
}

func (client *Client) CreateForwardLBFourthLayerListeners(args *CreateForwardLBFourthLayerListenersArgs) (
	*CreateForwardLBFourthLayerListenersResponse,
	error,
) {
	response := &CreateForwardLBFourthLayerListenersResponse{}
	err := client.Invoke("CreateForwardLBFourthLayerListeners", args, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type DeleteForwardLBListenerArgs struct {
	LoadBalancerId string `qcloud_arg:"loadBalancerId"`
	ListenerId     string `qcloud_arg:"listenerId"`
}

type DeleteForwardLBListenerResponse struct {
	Response
	RequestId int `json:"requestId"`
}

func (response DeleteForwardLBListenerResponse) Id() int {
	return response.RequestId
}

func (client *Client) DeleteForwardLBListener(args *DeleteForwardLBListenerArgs) (
	*DeleteForwardLBListenerResponse,
	error,
) {
	response := &DeleteForwardLBListenerResponse{}
	err := client.Invoke("DeleteForwardLBListener", args, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type DescribeForwardLBListenersArgs struct {
	LoadBalancerId   string    `qcloud_arg:"loadBalancerId"`
	ListenerIds      *[]string `qcloud_arg:"listenerIds"`
	Protocol         *int      `qcloud_arg:"protocol"`
	LoadBalancerPort *int      `qcloud_arg:"loadBalancerPort"`
}

type DescribeForwardLBListenersResponse struct {
	Response
	ListenerSet []FourthOrSeventeenLayerListener `json:"listenerSet"`
}

type FourthOrSeventeenLayerListener struct {
	ListenerId       string `json:"listenerid"`
	Protocol         int    `json:"protocol"`
	ProtocolType     string `json:"protocoltype"`
	LoadBalancerPort int    `json:"loadbalancerport"`
	SSLMode          string `json:"sslmode"`
	CertId           string `json:"certid"`
	CertCaId         string `json:"certcaid"`
}

func (client *Client) DescribeForwardLBListeners(args *DescribeForwardLBListenersArgs) (
	*DescribeForwardLBListenersResponse,
	error,
) {
	response := &DescribeForwardLBListenersResponse{}
	err := client.Invoke("DescribeForwardLBListeners", args, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type ModifyForwardFourthBackendsPortArgs struct {
	LoadBalancerId string               `qcloud_arg:"loadBalancerId"`
	ListenerId     string               `qcloud_arg:"listenerId"`
	Backends       []FourthBackendsPort `qcloud_arg:"backends"`
}

type FourthBackendsPort struct {
	InstanceId string `qcloud_arg:"instanceId"`
	Port       int    `qcloud_arg:"port"`
	NewPort    int    `qcloud_arg:"newPort"`
	Weight     *int   `qcloud_arg:"weight"`
}

type ModifyForwardFourthBackendsPortResponse struct {
	Response
	RequestId int `json:"requestId"`
}

func (client *Client) ModifyForwardFourthBackendsPort(args *ModifyForwardFourthBackendsPortArgs) (
	*ModifyForwardFourthBackendsPortResponse,
	error,
) {
	response := &ModifyForwardFourthBackendsPortResponse{}
	err := client.Invoke("ModifyForwardFourthBackendsPort", args, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type DescribeForwardLBBackendsArgs struct {
	LoadBalancerId   string    `qcloud_arg:"loadBalancerId"`
	ListenerIds      *[]string `qcloud_arg:"listenerIds"`
	Protocol         *int      `qcloud_arg:"protocol"`
	LoadBalancerPort *int      `qcloud_arg:"loadBalancerPort"`
}

type DescribeForwardLBBackendsResponse struct {
	Response
	Data []ForwardLBListener `json:"data"`
}

type ForwardLBListener struct {
	ListenerId       string                     `json:"listenerId"`
	Protocol         int                        `json:"protocol"`
	ProtocolType     string                     `json:"protocolType"`
	LoadBalancerPort int                        `json:"loadBalancerPort"`
	Rules            []ForwardLBListenerRule    `json:"rules"`
	Backends         []ForwardLBListenerBackend `json:"backends"`
}

type ForwardLBListenerRule struct {
	LocationId string                     `json:"locationId"`
	Domain     string                     `json:"domain"`
	Url        string                     `json:"url"`
	Backends   []ForwardLBListenerBackend `json:"backends"`
}

type ForwardLBListenerBackend struct {
	LanIp          string   `json:"LanIp"`
	WanIpSet       []string `json:"WanIpSet"`
	Port           int      `json:"Port"`
	Weight         int      `json:"Weight"`
	InstanceStatus int      `json:"InstanceStatus"`
	UnInstanceId   string   `json:"UnInstanceId"`
	InstanceName   string   `json:"InstanceName"`
	AddTimestamp   string   `json:"AddTimestamp"`
}

func (client *Client) DescribeForwardLBBackends(args *DescribeForwardLBBackendsArgs) (
	*DescribeForwardLBBackendsResponse,
	error,
) {
	response := &DescribeForwardLBBackendsResponse{}
	err := client.Invoke("DescribeForwardLBBackends", args, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type RegisterInstancesWithForwardLBFourthListenerArgs struct {
	LoadBalancerId string                                                    `qcloud_arg:"loadBalancerId"`
	ListenerId     string                                                    `qcloud_arg:"listenerId"`
	Backends       []RegisterInstancesWithForwardLBFourthListenerBackendOpts `qcloud_arg:"backends"`
}

type RegisterInstancesWithForwardLBFourthListenerBackendOpts struct {
	InstanceId string `qcloud_arg:"instanceId"`
	Port       int    `qcloud_arg:"port"`
	Weight     *int   `qcloud_arg:"weight"`
}

type RegisterInstancesWithForwardLBFourthListenerResponse struct {
	Response
	RequestId int `json:"requestId"`
}

func (response RegisterInstancesWithForwardLBFourthListenerResponse) Id() int {
	return response.RequestId
}

func (client *Client) RegisterInstancesWithForwardLBFourthListener(args *RegisterInstancesWithForwardLBFourthListenerArgs) (
	*RegisterInstancesWithForwardLBFourthListenerResponse,
	error,
) {
	response := &RegisterInstancesWithForwardLBFourthListenerResponse{}
	err := client.Invoke("RegisterInstancesWithForwardLBFourthListener", args, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type DeregisterInstancesFromForwardLBFourthListenerArgs struct {
	LoadBalancerId string                                                      `qcloud_arg:"loadBalancerId"`
	ListenerId     string                                                      `qcloud_arg:"listenerId"`
	Backends       []DeregisterInstancesWithForwardLBFourthListenerBackendOpts `qcloud_arg:"backends"`
}

type DeregisterInstancesWithForwardLBFourthListenerBackendOpts struct {
	InstanceId string `qcloud_arg:"instanceId"`
	Port       int    `qcloud_arg:"port"`
	Weight     *int   `qcloud_arg:"weight"`
}

type DeregisterInstancesFromForwardLBFourthListenerResponse struct {
	Response
	RequestId int `json:"requestId"`
}

func (response *DeregisterInstancesFromForwardLBFourthListenerResponse) Id() int {
	return response.RequestId
}

func (client *Client) DeregisterInstancesFromForwardLBFourthListener(args *DeregisterInstancesFromForwardLBFourthListenerArgs) (
	*DeregisterInstancesFromForwardLBFourthListenerResponse,
	error,
) {
	response := &DeregisterInstancesFromForwardLBFourthListenerResponse{}
	err := client.Invoke("DeregisterInstancesFromForwardLBFourthListener", args, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
