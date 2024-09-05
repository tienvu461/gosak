/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type FormatConfig struct {
	Json bool
	Text string
}

func JsonFormat(t string) error {
	fmt.Println("Output:\n", t)

	return nil
}
func FormatCreateCommand(c *FormatConfig) *cobra.Command {
	command := &cobra.Command{
		Use:   "format",
		Short: "format input (json & dict)",
		Long: `Take unformatted input (json & dict), 
		determine the type, remove escape charater, then spit out formatted text`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t := args[0]
			c.Text = t
			if c.Json {
				return JsonFormat(c.Text)
			} else {
				return nil
			}
		},
	}

	return command
}

func init() {
	config := &FormatConfig{
		Json: true,
		Text: "",
	}
	formatCmd := FormatCreateCommand(config)
	rootCmd.AddCommand(formatCmd)
}
