package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tencentcloud/tencentcloud-cloud-controller-manager/route-ctl/cmd/route"
	"github.com/tencentcloud/tencentcloud-cloud-controller-manager/route-ctl/cmd/routetable"
)

var rootCmd = &cobra.Command{
	Use:   "route-ctl",
	Short: "route command line tool to create vpc route table in tencentcloud",
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

func init() {
	routeTableCmd := &cobra.Command{
		Use:   "route-table",
		Short: "route table releated operations",
	}

	routeTableCmd.AddCommand(routetable.ListCmd)
	routeTableCmd.AddCommand(routetable.CreateCmd)
	routeTableCmd.AddCommand(routetable.DeleteCmd)

	routeCmd := &cobra.Command{
		Use:   "route",
		Short: "route releated operations",
	}

	routeCmd.AddCommand(route.ListCmd)
	routeCmd.AddCommand(route.CreateCmd)
	routeCmd.AddCommand(route.DeleteCmd)

	rootCmd.AddCommand(routeTableCmd)
	rootCmd.AddCommand(routeCmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
