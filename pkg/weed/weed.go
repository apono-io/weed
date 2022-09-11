package weed

import "fmt"

type Weed struct {
	Missing     []string
	Unnecessary []string
}

type Client struct {
	AwsService AwsService
}

func New(profile string) (client Client, err error) {
	service, err := NewAwsService(profile)
	if err != nil {
		return client, fmt.Errorf("error creating aws service: %v", err)
	}

	return Client{
		AwsService: service,
	}, err
}

func (c *Client) Find(desiredPermissions []string, iamRole string) (weed Weed, err error) {
	rolePermissions, err := c.AwsService.RolePermissions(iamRole)
	if err != nil {
		return weed, fmt.Errorf("error getting role actions: %v", err)
	}

	weed.Missing = diff(desiredPermissions, rolePermissions)
	weed.Unnecessary = diff(rolePermissions, desiredPermissions)

	return
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
