/*
Copyright © 2025 tienvu461@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"github.com/tienvu461/gosak/utils"
)

// ec2Cmd represents the ec2 command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Manage ec2 instance",
	Long:  `Manage ec2 instance with list, start, stop command`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List ec2 instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := getEc2Client()
		if err != nil {
			return err
		}

		input := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("instance-state-name"),
					Values: []*string{
						aws.String("pending"),
						aws.String("running"),
						aws.String("stopping"),
						aws.String("stopped"),
						aws.String("shutting-down"),
						aws.String("terminated"),
					},
				},
			},
		}
		result, err := svc.DescribeInstances(input)
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tState\tPrivateIP\tPublicIP")
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				name := ""
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						name = *tag.Value
						break
					}
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					*instance.InstanceId,
					name,
					*instance.State.Name,
					aws.StringValue(instance.PrivateIpAddress),
					aws.StringValue(instance.PublicIpAddress),
				)
			}
		}
		w.Flush()
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop <instance-id>|<instance-name>",
	Short: "Stop an ec2 instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := getEc2Client()
		if err != nil {
			return err
		}

		instanceID, err := resolveSingleInstanceID(svc, args[0], "running")
		if err != nil {
			return err
		}

		input := &ec2.StopInstancesInput{
			InstanceIds: []*string{
				aws.String(instanceID),
			},
		}

		_, err = svc.StopInstances(input)
		if err != nil {
			return err
		}
		fmt.Printf("Stopping instance %s\n", instanceID)
		return nil
	},
}

var startCmd = &cobra.Command{
	Use:   "start <instance-id>|<instance-name>",
	Short: "Start an ec2 instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := getEc2Client()
		if err != nil {
			return err
		}

		instanceID, err := resolveSingleInstanceID(svc, args[0], "stopped")
		if err != nil {
			return err
		}

		input := &ec2.StartInstancesInput{
			InstanceIds: []*string{
				aws.String(instanceID),
			},
		}

		_, err = svc.StartInstances(input)
		if err != nil {
			return err
		}
		fmt.Printf("Starting instance %s\n", instanceID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
	ec2Cmd.AddCommand(listCmd)
	ec2Cmd.AddCommand(stopCmd)
	ec2Cmd.AddCommand(startCmd)
}

func getEc2Client() (*ec2.EC2, error) {
	sess, err := utils.GetAWSSession()
	if err != nil {
		return nil, err
	}

	return ec2.New(sess), nil
}

func resolveSingleInstanceID(svc *ec2.EC2, identifier string, expectedStates ...string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{},
	}

	if len(expectedStates) > 0 {
		values := []*string{}
		for _, s := range expectedStates {
			values = append(values, aws.String(s))
		}
		input.Filters = append(input.Filters, &ec2.Filter{
			Name:   aws.String("instance-state-name"),
			Values: values,
		})
	}

	if strings.HasPrefix(identifier, "i-") {
		input.InstanceIds = []*string{aws.String(identifier)}
	} else {
		input.Filters = append(input.Filters, &ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: []*string{aws.String(identifier)},
		})
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return "", err
	}

	var foundInstances []*ec2.Instance
	for _, r := range result.Reservations {
		// Filter out terminated instances to avoid confusion
		for _, i := range r.Instances {
			if *i.State.Name != "terminated" {
				foundInstances = append(foundInstances, i)
			}
		}
	}

	if len(foundInstances) == 0 {
		return "", fmt.Errorf("no instance found with identifier %s matching criteria (states: %v)", identifier, expectedStates)
	}
	if len(foundInstances) > 1 {
		return "", fmt.Errorf("multiple active instances found with name %s, please use Instance ID instead", identifier)
	}

	return *foundInstances[0].InstanceId, nil
}
