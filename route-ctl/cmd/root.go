package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "route-ctl",
	Short: "route command line tool to create vpc route table in tencentcloud",
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
