/*
Copyright © 2022 tienvu461@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tienvu461/gosak/utils"
)

// versionCmd represents the version command
// To verify version and build info
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  `Show version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:\t", utils.Version)
		fmt.Println("Hash:\t\t", utils.Hash)
		fmt.Println("Operating System:\t", utils.Os)
		fmt.Println("Architect:\t", utils.Arch)
		fmt.Println("Golang Version:\t", utils.GoVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
