package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/dbdd4us/qcloudapi-sdk-go/ccs"
	"github.com/dbdd4us/qcloudapi-sdk-go/common"
	"github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&routeTableName, "route-table-name", "", "name of the route table to create")
	createCmd.Flags().StringVar(&routeTableCidrBlock, "route-table-cidr-block", "", "cidr of the route table to create")
	createCmd.Flags().StringVar(&vpcId, "vpc-id", "", "vpc id of the route table to create")
	createCmd.Flags().BoolVar(&ignoreConflict, "ignore-cidr-conflict", false, "Default false, ignore any cidr conflict when create route table. (Caution: This may results in disaster, use it very carefully)")

	createCmd.MarkFlagRequired("route-table-name")
	createCmd.MarkFlagRequired("route-table-cidr-block")
	createCmd.MarkFlagRequired("vpc-id")

}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create route table",
	RunE: func(cmd *cobra.Command, args []string) error {
		secretId := os.Getenv("QCloudSecretId")
		secretKey := os.Getenv("QCloudSecretKey")
		region := os.Getenv("QCloudCcsAPIRegion")
		logger := logrus.New()
		//logger.SetLevel(logrus.ErrorLevel)

		client, err := ccs.NewClient(common.Credential{SecretId: secretId, SecretKey: secretKey}, common.Opts{Logger: logger, Region: region})
		if err != nil {
			return err
		}

		ignoreConflictOpt := 0
		if ignoreConflict {
			ignoreConflictOpt = 1
		}
		_, err = client.CreateClusterRouteTable(&ccs.CreateClusterRouteTableArgs{
			RouteTableName:            routeTableName,
			RouteTableCidrBlock:       routeTableCidrBlock,
			VpcId:                     vpcId,
			IgnoreClusterCIDRConflict: ignoreConflictOpt,
		})
		return err
	},
}

var (
	routeTableName      string
	routeTableCidrBlock string
	vpcId               string
	ignoreConflict      bool
)
