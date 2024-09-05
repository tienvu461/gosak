/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

type AssumeConfig struct {
	Role     string
	Duration time.Duration
	Output   bool
}

func AssumeRole(role string, dura time.Duration) (*credentials.Value, error) {
	sess := session.Must(session.NewSession())

	svc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),
		RoleSessionName: aws.String("cli"),
		DurationSeconds: aws.Int64(int64(dura)),
	}

	resp, err := svc.AssumeRole(params)

	if err != nil {
		return nil, err
	}

	var creds credentials.Value
	creds.AccessKeyID = *resp.Credentials.AccessKeyId
	creds.SecretAccessKey = *resp.Credentials.SecretAccessKey
	creds.SessionToken = *resp.Credentials.SessionToken

	return &creds, nil
}

// assumeCmd represents the assume command
func AssumeCreateCommand(c *AssumeConfig) *cobra.Command {
	command := &cobra.Command{
		Use:   "assume",
		Short: "Execute sts assume and update/print token",
		Long:  `Execute sts assume and update/print token`,
		// Args:  AssumeArgsValidator(c),
		RunE: func(cmd *cobra.Command, args []string) error {
			var creds *credentials.Value
			creds, err := AssumeRole(c.Role, c.Duration)
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}
			fmt.Printf("export AWS_ACCESS_KEY_ID=\"%s\"\n", creds.AccessKeyID)
			fmt.Printf("export AWS_SECRET_ACCESS_KEY=\"%s\"\n", creds.SecretAccessKey)
			fmt.Printf("export AWS_SESSION_TOKEN=\"%s\"\n", creds.SessionToken)
			fmt.Printf("export AWS_SECURITY_TOKEN=\"%s\"\n", creds.SessionToken)
			fmt.Printf("export ASSUMED_ROLE=\"%s\"\n", c.Role)
			fmt.Printf("# Run this to configure your shell:\n")
			fmt.Printf("# eval $(%s)\n", strings.Join(os.Args, " "))
			return nil
		},
	}

	return command
}

func init() {
	config := &AssumeConfig{
		Role:     "",
		Duration: 3600,
		Output:   false,
	}
	assumeCmd := AssumeCreateCommand(config)
	rootCmd.AddCommand(assumeCmd)
}
