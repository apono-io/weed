/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/apono-io/weed/pkg/weed"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var permissions []string
var roleArn string
var failOnDiff bool
var failOnMissing bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weed",
	Short: "A brief description of weed application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := findWeed(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func findWeed() (err error) {
	weedClient, err := weed.New()
	if err != nil {
		return
	}

	diff, err := weedClient.Find(permissions, roleArn)
	if err != nil {
		return
	}

	if len(diff.Added) > 0 {
		fmt.Printf("Added %d permissions: \n", len(diff.Added))
		for _, perm := range diff.Added {
			color.Green(fmt.Sprintf("  %s", perm))
		}
	}

	if len(diff.Removed) > 0 {
		fmt.Printf("Removed %d permissions: \n", len(diff.Removed))
		for _, perm := range diff.Removed {
			color.Red(fmt.Sprintf("  %s", perm))
		}
	}

	if len(diff.Added) > 0 {
		if failOnMissing {
			os.Exit(1)
		}

		if len(diff.Removed) > 0 && failOnDiff {
			os.Exit(1)
		}
	}

	return
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringSliceVarP(&permissions, "permissions", "p", []string{}, "Desired permissions")
	rootCmd.Flags().StringVarP(&roleArn, "role-arn", "r", "", "Role ARN")
	rootCmd.Flags().BoolVarP(&failOnDiff, "fail-on-diff", "f", false, "Return error if diff is found")
	rootCmd.Flags().BoolVarP(&failOnMissing, "fail-on-missing", "m", false, "Return error if permissions are missing")
}
