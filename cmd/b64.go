/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tienvu461/gosak/utils"
)

type B64Config struct {
	Encode bool
	Text   string
}

func B64Decode(t string) error {
	fmt.Println("Decoding:", t)
	d_text, err := base64.StdEncoding.DecodeString(t)

	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Unable to decode", t)
		// return utils.Code(1)
		return err
	}
	fmt.Printf("%q\n", d_text)
	return err
}
func B64Encode(t string) error {
	fmt.Println("Encoding:", t)
	e_text := base64.StdEncoding.EncodeToString([]byte(t))

	fmt.Println(e_text)
	return nil
}
func B64CreateCommand(c *B64Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "b64",
		Short: "Base64 encode, decode",
		Long:  `Base64 encode, decode`,
		Args:  B64ArgsValidator(c),
		RunE: func(cmd *cobra.Command, args []string) error {
			// fmt.Println(c.Text)
			// str := base64.StdEncoding.EncodeToString([]byte("Hello, playground"))
			// fmt.Println(str)
			if c.Encode {
				return B64Encode(c.Text)
			} else {
				return B64Decode(c.Text)
			}
		},
	}

	return command
}

func init() {
	config := &B64Config{
		Encode: false,
		Text:   "abc",
	}
	b64Cmd := B64CreateCommand(config)
	rootCmd.AddCommand(b64Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// b64Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// b64Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	b64Cmd.Flags().BoolVarP(&config.Encode, "encode", "e", false, "Encode base64")
	// b64Cmd.Flags().StringVarP(&config.Text, "text", "t", "", "Text")
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
