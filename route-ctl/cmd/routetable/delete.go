package routetable

import (
	"os"

	"github.com/dbdd4us/qcloudapi-sdk-go/ccs"
	"github.com/dbdd4us/qcloudapi-sdk-go/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	DeleteCmd.Flags().StringVar(&routeTableNameToDelete, "route-table-name", "", "name of the route table to create")

	DeleteCmd.MarkFlagRequired("route-table-name")

}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete route table",
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
		_, err = client.DeleteClusterRouteTable(&ccs.DeleteClusterRouteTableArgs{
			RouteTableName: routeTableNameToDelete,
		})
		return err
	},
}

var (
	routeTableNameToDelete string
)
