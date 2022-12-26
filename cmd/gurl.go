/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
    "os"
    "net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
    "github.com/tienvu461/gosak/gurl"
    "github.com/tienvu461/gosak/utils"
)

func CreateCommand(c *gurl.Config, headers []string) *cobra.Command {
    command := &cobra.Command{
        Use:   "gurl <URL>",
        Short: "Cheap knockoff of famous curl",
        Long: `gURL is the first command client which similar to cURL HTTP Client that I built following this tutorial https://dev.to/mauriciolinhares/building-and-distributing-a-command-line-tool-in-golang-go0
            Usage:
        `,
        Args:    ArgsValidator(c),
        PreRunE: OptionsValidator(c, headers),
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Println("gurl called")
            return gurl.Execute(c)
        },
    }

   return command
}
func init() {
    headers := make([]string, 0, 255)
    config := &gurl.Config {
          Headers:            map[string][]string{},
          ResponseBodyOutput: os.Stdout,
          ControlOutput:      os.Stdout,
    }
    gurlCmd := CreateCommand(config, headers)
    rootCmd.AddCommand(gurlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gurlCmd.PersistentFlags().String("foo", "", "A help for foo")

    gurlCmd.PersistentFlags().StringSliceVarP(&headers, "headers", "H", nil, `custom headers headers to be sent with the request, headers are separated by "," as in "HeaderName: Header content,OtherHeader: Some other value"`)
    gurlCmd.PersistentFlags().StringVarP(&config.UserAgent, "user-agent", "u", "gurl", "the user agent to be used for requests")
    gurlCmd.PersistentFlags().StringVarP(&config.Data, "data", "d", "", "data to be sent as the request body")
    gurlCmd.PersistentFlags().StringVarP(&config.Method, "method", "m", http.MethodGet, "HTTP method to be used for the request")
    gurlCmd.PersistentFlags().BoolVarP(&config.Insecure, "insecure", "k", false, "allows insecure server connections over HTTPS")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gurlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ArgsValidator(c *gurl.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if l := len(args); l != 1 {
			return utils.NewErrorWithCode(2, "you must provide a single URL to be called but you provided %v", l)
		}

		u, err := url.Parse(args[0])
		if err != nil {
			return errors.Wrapf(err, "the URL provided is invalid: %v", args[0])
		}

		c.Url = u

		return nil
	}
}

func OptionsValidator(c *gurl.Config, headers []string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, h := range headers {
			if name, value, found := strings.Cut(h, ":"); found {
				c.Headers.Add(strings.TrimSpace(name), strings.TrimSpace(value))
			} else {
				return utils.NewErrorWithCode(3, "header is not a valid http header separated by `:`, value was: [%v]", h)
			}
		}

		return nil
	}
}
