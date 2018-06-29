package cmd

import (
	"fmt"
	"net"
	"os"
	"text/tabwriter"

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

		_, cidr, err := net.ParseCIDR(routeTableCidrBlock)
		if err != nil {
			return fmt.Errorf("invalid cidr %s", routeTableCidrBlock)
		}

		client, err := ccs.NewClient(common.Credential{SecretId: secretId, SecretKey: secretKey}, common.Opts{Logger: logger, Region: region})
		if err != nil {
			return err
		}

		response, err := client.CheckClusterRouteTableCidrConflict(&ccs.CheckClusterRouteTableCidrConflictArgs{
			RouteTableCidrBlock: cidr.String(),
			VpcId:               vpcId,
		})
		if err != nil {
			return err
		}

		if response.Data.HasConflict {
			logger.Errorf("Cidr %s has following conflicts. Use --ignore-cidr-conflict to ignore conflict. (Caution: This operation is dangerous and may results in disaster, use it very carefully)\n", routeTableCidrBlock)
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', tabwriter.TabIndent)
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t", "ConflictType", "Name", "Id", "Cidr"))
			for _, conflict := range response.Data.CidrConflicts {
				fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t", conflict.Type, conflict.Name, conflict.Id, conflict.Cidr))
			}
			w.Flush()

			if !ignoreConflict {
				return nil
			}
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
