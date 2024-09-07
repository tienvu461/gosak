/*
Copyright © 2022 tienvu461@gmail.com
*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tienvu461/gosak/utils"
)

// sslcertCmd represents the sslcert command
var sslcertCmd = &cobra.Command{
	Use:   "sslcert <domain>:443",
	Short: "Get SSL certificate status, expire date, TLS version, ...",
	Long:  `Get SSL certificate status, expire date, TLS version, ...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if l := len(args); l != 1 {
			return utils.NewErrorWithCode(2, "you must provide a single URL to be called but you provided %v", l)
		}

		u, err := url.Parse(args[0])
		if err != nil {
			return errors.Wrapf(err, "the URL provided is invalid: %v", args[0])
		}
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}
		//TODO: automatically parsing domain name from any input type and
		//appending port 443 afterward
		conn, err := tls.Dial("tcp", u.String(), conf)
		if err != nil {
			fmt.Println("Error in Dial", err)
			return err
		}
		defer conn.Close()
		certs := conn.ConnectionState().PeerCertificates
		for _, cert := range certs {
			fmt.Printf("Issuer Name: %s\n", cert.Issuer)
			fmt.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-January-02"))
			fmt.Printf("Common Name: %s \n", cert.Issuer.CommonName)

		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sslcertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sslcertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sslcertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
