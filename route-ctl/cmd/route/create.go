package route

import (
	"errors"
	"net"
	"os"

	"fmt"
	"github.com/dbdd4us/qcloudapi-sdk-go/ccs"
	"github.com/dbdd4us/qcloudapi-sdk-go/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	CreateCmd.Flags().StringVar(&routeTableName, "route-table-name", "", "name of the route table to create route in")
	CreateCmd.Flags().StringVar(&destinationCidrBlock, "destination-cidr-block", "", "route destination cidr block")
	CreateCmd.Flags().StringVar(&gatewayIp, "gateway-ip", "", "route gateway ip")

	CreateCmd.MarkFlagRequired("route-table-name")

}

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create route",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, _, err := net.ParseCIDR(destinationCidrBlock); err != nil {
			return errors.New(fmt.Sprintf("Invalid destination cidr block %s, err: %s", destinationCidrBlock, err))
		}

		secretId := os.Getenv("QCloudSecretId")
		secretKey := os.Getenv("QCloudSecretKey")
		region := os.Getenv("QCloudCcsAPIRegion")
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)

		client, err := ccs.NewClient(common.Credential{SecretId: secretId, SecretKey: secretKey}, common.Opts{Logger: logger, Region: region})
		if err != nil {
			return err
		}

		rtbResponse, err := client.DescribeClusterRouteTable(&ccs.DescribeClusterRouteTableArgs{})
		if err != nil {
			return err
		}

		routeTable := new(ccs.RouteTableInfo)

		for _, rtb := range rtbResponse.Data.RouteTableSet {
			if rtb.RouteTableName == routeTableName {
				routeTable = &rtb
				break
			}
		}

		if routeTable == nil {
			return errors.New(fmt.Sprintf("route table %s not found", routeTableName))
		}

		_, err = client.CreateClusterRoute(&ccs.CreateClusterRouteArgs{
			RouteTableName:       routeTableName,
			DestinationCidrBlock: destinationCidrBlock,
			GatewayIp:            gatewayIp,
		})
		return err
	},
}

var (
	routeTableName       string
	destinationCidrBlock string
	gatewayIp            string
)
