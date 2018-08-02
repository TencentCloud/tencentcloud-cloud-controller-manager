package route

import (
	"os"

	"fmt"
	"github.com/dbdd4us/qcloudapi-sdk-go/ccs"
	"github.com/dbdd4us/qcloudapi-sdk-go/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"text/tabwriter"
)

func init() {
	ListCmd.Flags().StringVar(&routeTableName, "route-table-name", "", "name of the route table to create route in")

	ListCmd.MarkFlagRequired("route-table-name")

}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list route",
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
		response, err := client.DescribeClusterRoute(&ccs.DescribeClusterRouteArgs{
			RouteTableName: routeTableName,
		})
		if err != nil {
			return err
		}

		if len(response.Data.RouteSet) == 0 {
			fmt.Println("No route found")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t", "RouteTableName", "DestinationCidrBlock", "GatewayIp"))
		for _, route := range response.Data.RouteSet {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t", route.RouteTableName, route.DestinationCidrBlock, route.GatewayIp))
		}
		w.Flush()
		return nil
	},
}
