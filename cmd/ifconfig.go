/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/tienvu461/gosak/gurl"
	"github.com/tienvu461/gosak/utils"
)

// ifconfigCmd represents the ifconfig command
var ifconfigCmd = &cobra.Command{
	Use:   "ifconfig",
	Short: "Get current Public IP",
	Long:  `Utilize gurl cmd to make a HTTP Get Request to https://ifconfig.me and return a Public IP`,
	RunE: func(cmd *cobra.Command, args []string) error {
		u, err := url.Parse(utils.IFCONFIG_URL)
		if err != nil {
			return nil
		}
		config := &gurl.Config{
			Headers:            map[string][]string{},
			UserAgent:          "curl",
			ResponseBodyOutput: os.Stdout,
			ControlOutput:      os.Stdout,
			Url:                u,
		}
		return gurl.Execute(config)
	},
}

func init() {
	rootCmd.AddCommand(ifconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ifconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ifconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
