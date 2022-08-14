package cmd

import (
	"fmt"
	"github.com/apono-io/weed/pkg/build"
	"github.com/spf13/cobra"
)

var VersionCommand = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("WEED version: %s (%s) build date: %s\n", build.Version, build.Commit, build.Date)
	},
}
