package cmd

import (
	"text/tabwriter"
	"os"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/dbdd4us/qcloudapi-sdk-go/ccs"
	"github.com/dbdd4us/qcloudapi-sdk-go/common"
	"github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all route table",
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
		response, err := client.DescribeClusterRouteTable(&ccs.DescribeClusterRouteTableArgs{})
		if err != nil {
			return err
		}
		if len(response.Data.RouteTableSet) == 0 {
			fmt.Println("No route table found")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t", "RouteTableName", "RouteTableCidrBlock", "VpcId"))
		for _, routeTable := range response.Data.RouteTableSet {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t", routeTable.RouteTableName, routeTable.RouteTableCidrBlock, routeTable.VpcId))
		}
		w.Flush()
		return nil
	},
}
