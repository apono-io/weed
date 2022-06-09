package weed

import "fmt"

type Weed struct {
	Added   []string
	Removed []string
}

type WeedClient struct {
	AwsService AwsService
}

func New() (client WeedClient, err error) {
	service, err := NewAwsService()
	if err != nil {
		return client, fmt.Errorf("error creating aws service: %v", err)
	}

	return WeedClient{
		AwsService: service,
	}, err
}

func (c *WeedClient) Find(desiredPermissions []string, iamRole string) (weed Weed, err error) {
	rolePermissions, err := c.AwsService.RolePermissions(iamRole)
	if err != nil {
		return weed, fmt.Errorf("error getting role permissions: %v", err)
	}

	weed.Added = diff(desiredPermissions, rolePermissions)
	weed.Removed = diff(rolePermissions, desiredPermissions)

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
