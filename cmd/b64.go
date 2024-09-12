/*
Copyright © 2024 tienvu461@gmail.com
*/
package cmd

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tienvu461/gosak/utils"
)

type B64Config struct {
	Encode bool
	Text   string
}

func B64Decode(t string) (string, error) {
	fmt.Println("Decoding:", t)
	d_text, err := base64.StdEncoding.DecodeString(t)
	d_str := strings.ReplaceAll(string(d_text[:]), "\n", "")

	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Unable to decode", t)
		// return utils.Code(1)
		return "", err
	}
	fmt.Printf("%q\n", d_str)
	return d_str, nil
}
func B64Encode(t string) (string, error) {
	fmt.Println("Encoding:", t)
	e_text := base64.StdEncoding.EncodeToString([]byte(t))

	fmt.Println(e_text)
	return e_text, nil
}
func B64CreateCommand(c *B64Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "b64",
		Short: "Base64 encode, decode",
		Long:  `Take text or b64 encoded input and do encode/decode`,
		Args:  B64ArgsValidator(c),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var res string
			if c.Encode {
				res, err = B64Encode(c.Text)
			} else {
				res, err = B64Decode(c.Text)
			}
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
		},
	}

	return command
}

func init() {
	config := &B64Config{
		Encode: false,
		Text:   "",
	}
	b64Cmd := B64CreateCommand(config)
	rootCmd.AddCommand(b64Cmd)
	// Flgas
	b64Cmd.Flags().BoolVarP(&config.Encode, "encode", "e", false, "Encode base64")
}

func B64ArgsValidator(c *B64Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if l := len(args); l != 1 {
			return utils.NewErrorWithCode(2, "you must provide a text but you provided %v", l)
		}
		t := args[0]
		if !c.Encode && len(t)%4 != 0 {
			return utils.NewErrorWithCode(2, "you must provide a valid base64 encoded text but you provided %v", t)
		}
		c.Text = t
		return nil
	}
}
