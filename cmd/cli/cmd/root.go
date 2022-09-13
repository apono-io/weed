package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/apono-io/weed/pkg/weed"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func Execute() error {
	var actions []string
	var awsProfile string
	var roleArn string
	var policyFile string
	var remoteRoleArn string
	var remoteRoleAwsProfile string
	var failOnDiff bool
	var failOnMissing bool

	var rootCmd = &cobra.Command{
		Use: "weed",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(actions) == 0 && policyFile == "" && remoteRoleArn == "" {
				return errors.New("required flag(s) \"actions, policy-file, remote-role\" not set")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if policyFile != "" {
				file, err := os.Open(policyFile)
				if err != nil {
					return err
				}

				defer func(file *os.File) {
					_ = file.Close()
				}(file)

				var policy weed.PolicyPermission
				err = json.NewDecoder(file).Decode(&policy)
				if err != nil {
					return err
				}

				for _, stmt := range policy.Statement {
					actions = append(actions, stmt.Action...)
				}
			}

			if remoteRoleArn != "" {
				weedClient, err := weed.New(remoteRoleAwsProfile)
				if err != nil {
					return err
				}

				rolePermissions, err := weedClient.AwsService.RolePermissions(remoteRoleArn)
				if err != nil {
					return err
				}

				actions = append(actions, rolePermissions...)
			}

			weedClient, err := weed.New(awsProfile)
			if err != nil {
				return err
			}

			diff, err := weedClient.Find(actions, roleArn)
			if err != nil {
				return err
			}

			if len(diff.Missing) > 0 {
				fmt.Printf("Missing %d actions:\n", len(diff.Missing))
				for _, perm := range diff.Missing {
					color.Red(fmt.Sprintf("  %s", perm))
				}
			}

			if len(diff.Unnecessary) > 0 {
				fmt.Printf("Unnecessary %d actions:\n", len(diff.Unnecessary))
				for _, perm := range diff.Unnecessary {
					text := fmt.Sprintf("  %s", perm)
					if failOnDiff {
						color.Red(text)
					} else {
						color.Green(text)
					}
				}
			}

			if len(diff.Missing) > 0 && failOnMissing {
				os.Exit(1)
			}

			if len(diff.Unnecessary) > 0 && failOnDiff {
				os.Exit(1)
			}

			if len(diff.Missing) == 0 && len(diff.Unnecessary) == 0 {
				fmt.Printf("Role %s is in sync\n", roleArn)
			}

			return nil
		},
	}

	flags := rootCmd.Flags()
	flags.StringVarP(&awsProfile, "aws-profile", "p", "", "Role AWS Profile")
	flags.StringVarP(&roleArn, "role", "r", "", "Role ARN/Name")
	err := rootCmd.MarkFlagRequired("role")
	if err != nil {
		return err
	}

	flags.StringSliceVarP(&actions, "actions", "a", []string{}, "Desired actions")
	flags.StringVarP(&policyFile, "policy-file", "f", "", "Role ARN")
	flags.StringVarP(&remoteRoleArn, "remote-role", "", "", "Remote Role ARN/Name")
	flags.StringVarP(&remoteRoleAwsProfile, "remote-role-aws-profile", "", "", "Remote Role AWS Profile")
	flags.BoolVarP(&failOnDiff, "fail-on-diff", "d", false, "Return error if diff is found")
	flags.BoolVarP(&failOnMissing, "fail-on-missing", "m", false, "Return error if actions are missing")

	rootCmd.AddCommand(VersionCommand)
	return rootCmd.Execute()
}
