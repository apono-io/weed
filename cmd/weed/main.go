package main

import (
	"fmt"
	"github.com/apono-io/weed/pkg/core"
)

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/davecgh/go-spew/spew"
)

func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

type ArrayOfStrings []string

// UnmarshalJSON implements the json.Unmarshaler interface
func (aos *ArrayOfStrings) UnmarshalJSON(b []byte) error {

	var actions []string

	if b[0] == '"' {
		var action string
		if err := json.Unmarshal(b, &action); err != nil {
			return err
		}
		actions = append(actions, action)
	} else if err := json.Unmarshal(b, &actions); err != nil {
		return err
	}

	*aos = actions

	return nil
}

type PolicyPermissionStatement struct {
	Effect   string         `json:"Effect"`
	Action   ArrayOfStrings `json:"Action"`
	Resource ArrayOfStrings `json:"Resource"`
	// Condition interface{}    `json:"Condition"` // TODO: add logic to handle this
}

type PolicyPermission struct {
	Version   string                      `json:"Version"`
	Statement []PolicyPermissionStatement `json:"Statement"`
}

type Diff struct {
	Added   []string
	Removed []string
}

func diff(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func main() {
	fmt.Println("CLI", core.Version, core.Commit, core.BuildDate)
}
