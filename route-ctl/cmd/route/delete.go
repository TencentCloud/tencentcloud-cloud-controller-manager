package route

import (
	"os"

	"github.com/dbdd4us/qcloudapi-sdk-go/ccs"
	"github.com/dbdd4us/qcloudapi-sdk-go/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	DeleteCmd.Flags().StringVar(&routeTableName, "route-table-name", "", "name of the route table to create route in")
	DeleteCmd.Flags().StringVar(&destinationCidrBlock, "destination-cidr-block", "", "route destination cidr block")
	DeleteCmd.Flags().StringVar(&gatewayIp, "gateway-ip", "", "route gateway ip")

	DeleteCmd.MarkFlagRequired("route-table-name")
	DeleteCmd.MarkFlagRequired("destination-cidr-block")
	DeleteCmd.MarkFlagRequired("gateway-ip")

}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete route",
	RunE: func(cmd *cobra.Command, args []string) error {
		secretId := os.Getenv("QCloudSecretId")
		secretKey := os.Getenv("QCloudSecretKey")
		region := os.Getenv("QCloudCcsAPIRegion")
		logger := logrus.New()
		logger.SetLevel(logrus.ErrorLevel)

		client, err := ccs.NewClient(common.Credential{SecretId: secretId, SecretKey: secretKey}, common.Opts{Logger: logger, Region: region})
		if err != nil {
			return err
		}
		_, err = client.DeleteClusterRoute(&ccs.DeleteClusterRouteArgs{
			RouteTableName:       routeTableName,
			DestinationCidrBlock: destinationCidrBlock,
			GatewayIp:            gatewayIp,
		})
		return err
	},
}
