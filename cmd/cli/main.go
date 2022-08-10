package main

import (
	"fmt"
	"github.com/apono-io/weed/cmd/cli/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
