/*
Copyright © 2024 tienvu461@gmail.com
Credit to https://github.com/remind101/assume-role/
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
	"github.com/tienvu461/gosak/utils"
)

var roleArnRe = regexp.MustCompile(`^arn:aws:iam::[0-9]*:role/[a-z0-9]*`)

const (
	AWS_STS_MAX_DURATION = 43200
)

type AssumeConfig struct {
	Role    string
	Profile string
	// Duration time.Duration
	Duration int64
	Output   bool
}

// Assume to given iam role
// using env or default credential from ~/.aws/credentials
func AssumeRole(role string, dura int64) (*credentials.Value, error) {
	sess := session.Must(session.NewSession())
	svc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),
		RoleSessionName: aws.String("gosak"),
		DurationSeconds: aws.Int64(dura),
	}

	resp, err := svc.AssumeRole(params)

	if err != nil {
		return nil, err
	}
	// fmt.Println(resp.Credentials.Expiration)
	var creds credentials.Value
	creds.AccessKeyID = *resp.Credentials.AccessKeyId
	creds.SecretAccessKey = *resp.Credentials.SecretAccessKey
	creds.SessionToken = *resp.Credentials.SessionToken

	return &creds, nil
}

// read mfa from terminal
func mfaReader() (string, error) {
    fmt.Fprintf(os.Stderr, "MFA token is required: ") // Printout message even eval() is wrapped
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1) // convert CRLF to LF
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// Assume to profile defined in ~/.aws/config
// using env or default credential from ~/.aws/credentials
func AssumeProfile(profile string, dura int64) (*credentials.Value, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:                 profile,
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: mfaReader,
	}))

	// fmt.Println("Profile Provider does not support getting Expiration")
	creds, err := sess.Config.Credentials.Get()

	if err != nil {
		return nil, err
	}
	return &creds, nil
}

// assumeCmd represents the assume command
func AssumeCreateCommand(c *AssumeConfig) *cobra.Command {
	command := &cobra.Command{
		Use:     "assume",
		Short:   "Execute sts assume and update/print token",
		Long:    `Execute sts assume and update/print token`,
		PreRunE: AssumeOptionsValidator(c),
		RunE: func(cmd *cobra.Command, args []string) error {
			var creds *credentials.Value
			var err error
			var assumeRoleName string
			// fmt.Println(c)
			if c.Role != "" { // Assume Role from given role arn
				creds, err = AssumeRole(c.Role, c.Duration)
				lastInd := strings.LastIndex(c.Role, "/")
				assumeRoleName = c.Role[lastInd+1:]
			} else if c.Profile != "" { // Assume Role from given profile name
				creds, err = AssumeProfile(c.Profile, c.Duration)
				assumeRoleName = c.Profile
			}
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}
			fmt.Printf("export AWS_ACCESS_KEY_ID=\"%s\"\n", creds.AccessKeyID)
			fmt.Printf("export AWS_SECRET_ACCESS_KEY=\"%s\"\n", creds.SecretAccessKey)
			fmt.Printf("export AWS_SESSION_TOKEN=\"%s\"\n", creds.SessionToken)
			fmt.Printf("export AWS_SECURITY_TOKEN=\"%s\"\n", creds.SessionToken)
			fmt.Printf("export ASSUMED_ROLE=\"%s\"\n", assumeRoleName)
			fmt.Printf("# Run this to configure your shell:\n")
			fmt.Printf("# eval $(%s)\n", strings.Join(os.Args, " "))
			return nil
		},
	}

	return command
}

func init() {
	config := &AssumeConfig{}
	// var duration int64
	assumeCmd := AssumeCreateCommand(config)
	rootCmd.AddCommand(assumeCmd)

	assumeCmd.Flags().StringVarP(&config.Role, "role-arn", "r", "", "Assume to given Role ARN. Either Role or Profile is specified")
	assumeCmd.Flags().StringVarP(&config.Profile, "profile", "p", "", "Assume to given Profile name. Either Role or Profile is specified")
	assumeCmd.Flags().Int64VarP(&config.Duration, "duration", "d", 3600, "Session duration")
	// assumeCmd.Flags().DurationVarP(&config.Duration, "duration", "d", 3600, "Session duration")
	assumeCmd.Flags().BoolVarP(&config.Output, "output", "o", false, "Output creds")
}

func AssumeOptionsValidator(c *AssumeConfig) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if c.Role != "" && c.Profile != "" {
			return utils.NewErrorWithCode(3, "Profile & Role ARN cannot be both specified")
		}
		if c.Role != "" && !roleArnRe.MatchString(c.Role) {
			return utils.NewErrorWithCode(3, "Role ARN is not in correct format: \"%v\"", c.Role)
		}
		if c.Duration > AWS_STS_MAX_DURATION {
			return utils.NewErrorWithCode(3, "Duration cannot exceed %v: \"%v\"", AWS_STS_MAX_DURATION, c.Duration)
		}
		return nil
	}
}
