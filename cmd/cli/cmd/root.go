package cmd

import (
	"fmt"
	"os"

	"github.com/apono-io/weed/pkg/weed"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func Execute() error {
	var permissions []string
	var roleArn string
	var failOnDiff bool
	var failOnMissing bool

	var rootCmd = &cobra.Command{
		Use: "weed",
		RunE: func(cmd *cobra.Command, args []string) error {
			weedClient, err := weed.New()
			if err != nil {
				return err
			}

			diff, err := weedClient.Find(permissions, roleArn)
			if err != nil {
				return err
			}

			if len(diff.Missing) > 0 {
				fmt.Printf("Missing %d permissions: \n", len(diff.Missing))
				for _, perm := range diff.Missing {
					color.Green(fmt.Sprintf("  %s", perm))
				}
			}

			if len(diff.Unnecessary) > 0 {
				fmt.Printf("Unnecessary %d permissions: \n", len(diff.Unnecessary))
				for _, perm := range diff.Unnecessary {
					color.Red(fmt.Sprintf("  %s", perm))
				}
			}

			if len(diff.Missing) > 0 {
				if failOnMissing {
					os.Exit(1)
				}

				if len(diff.Unnecessary) > 0 && failOnDiff {
					os.Exit(1)
				}
			}

			return nil
		},
	}

	flags := rootCmd.Flags()
	flags.StringVarP(&roleArn, "role-arn", "r", "", "Role ARN")
	err := rootCmd.MarkFlagRequired("role-arn")
	if err != nil {
		return err
	}

	flags.StringSliceVarP(&permissions, "permissions", "p", []string{}, "Desired permissions")
	flags.BoolVarP(&failOnDiff, "fail-on-diff", "f", false, "Return error if diff is found")
	flags.BoolVarP(&failOnMissing, "fail-on-missing", "m", false, "Return error if permissions are missing")

	return rootCmd.Execute()
}
