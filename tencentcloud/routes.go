package tencentcloud

import (
	"context"

	"github.com/dbdd4us/qcloudapi-sdk-go/ccs"

	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/apimachinery/pkg/types"
)

// ListRoutes lists all managed routes that belong to the specified clusterName
func (cloud *Cloud) ListRoutes(ctx context.Context, clusterName string) ([]*cloudprovider.Route, error) {
	cloudRoutes, err := cloud.ccs.DescribeClusterRoute(&ccs.DescribeClusterRouteArgs{RouteTableName: cloud.config.ClusterRouteTable})
	if err != nil {
		return []*cloudprovider.Route{}, err
	}

	routes := make([]*cloudprovider.Route, len(cloudRoutes.RouteSet))

	for idx, route := range cloudRoutes.RouteSet {
		routes[idx] = &cloudprovider.Route{Name: route.GatewayIp, TargetNode: types.NodeName(route.GatewayIp), DestinationCIDR: route.DestinationCidrBlock}
	}
	return routes, nil
}

// CreateRoute creates the described managed route
// route.Name will be ignored, although the cloud-provider may use nameHint
// to create a more user-meaningful name.
func (cloud *Cloud) CreateRoute(ctx context.Context, clusterName string, nameHint string, route *cloudprovider.Route) error {
	_, err := cloud.ccs.CreateClusterRoute(&ccs.CreateClusterRouteArgs{
		RouteTableName:       cloud.config.ClusterRouteTable,
		GatewayIp:            string(route.TargetNode),
		DestinationCidrBlock: route.DestinationCIDR,
	})

	return err
}

// DeleteRoute deletes the specified managed route
// Route should be as returned by ListRoutes
func (cloud *Cloud) DeleteRoute(ctx context.Context, clusterName string, route *cloudprovider.Route) error {
	_, err := cloud.ccs.DeleteClusterRoute(&ccs.DeleteClusterRouteArgs{
		RouteTableName:       route.Name,
		GatewayIp:            string(route.TargetNode),
		DestinationCidrBlock: route.DestinationCIDR,
	})
	return err
}
