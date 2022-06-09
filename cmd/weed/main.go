package main

import (
	"fmt"

	"github.com/apono-io/weed/pkg/core"
)

func main() {
	fmt.Println("CLI", core.Version, core.Commit, core.BuildDate)
}
