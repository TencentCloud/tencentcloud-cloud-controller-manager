package ccs

type CreateClusterRouteTableArgs struct {
	RouteTableName      string `qcloud_arg:"RouteTableName"`
	RouteTableCidrBlock string `qcloud_arg:"RouteTableCidrBlock"`
	VpcId               string `qcloud_arg:"VpcId"`
}

type CreateClusterRouteTableResponse struct {
	Response
}

type DeleteClusterRouteTableArgs struct {
	RouteTableName string `qcloud_arg:"RouteTableName"`
}

type DeleteClusterRouteTableResponse struct {
	Response
}

func (client *Client) CreateClusterRouteTable(args *CreateClusterRouteTableArgs) (*CreateClusterRouteTableResponse, error) {
	response := &CreateClusterRouteTableResponse{}
	err := client.Invoke("CreateClusterRouteTable", args, response)
	if err != nil {
		return &CreateClusterRouteTableResponse{}, err
	}
	return response, nil
}

func (client *Client) DeleteClusterRouteTable(args *DeleteClusterRouteTableArgs) (*DeleteClusterRouteTableResponse, error) {
	response := &DeleteClusterRouteTableResponse{}
	err := client.Invoke("DeleteClusterRouteTable", args, response)
	if err != nil {
		return &DeleteClusterRouteTableResponse{}, err
	}
	return response, nil
}
