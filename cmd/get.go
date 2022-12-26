/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
        var endPoint ="https://ifconfig.me"

        if len(args) >= 1 && args[0] != "" {
            endPoint = args[0]
        }

        fmt.Println("Trying to get public IP Address")

        response, err := http.Get(endPoint)

        if err != nil {
            fmt.Println(err)
            fmt.Println("Cannot get public IP Address")
        }

        defer response.Body.Close()

        if response.StatusCode == 200 {
            body, err := io.ReadAll(response.Body)
            if err != nil {
                fmt.Println(err)
            }
            fmt.Println(string(body))
        } else {
            fmt.Println("Fail to reach" + endPoint)
        }
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
