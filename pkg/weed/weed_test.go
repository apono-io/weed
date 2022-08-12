package weed

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestDiff(t *testing.T) {
	tests := []struct {
		current []string
		desired []string
		want    []string
	}{
		{current: []string{}, desired: []string{}, want: []string{}},
		{current: []string{"a", "b", "c"}, desired: []string{"a", "b", "c"}, want: []string{}},
		{current: []string{"a", "b", "c"}, desired: []string{"a", "b"}, want: []string{"c"}},
	}

	for _, test := range tests {
		got := diff(test.current, test.desired)
		if testEq(got, test.want) == false {
			t.Errorf("diff(%v, %v) = %v, want %v", test.current, test.desired, got, test.want)
		}
	}
}

func TestFindWeed(t *testing.T) {
	iamRole := "arn:aws:iam::123456789012:role/test-role"
	weedClient := &Client{
		AwsService: newMockAwsService(awsTestCases[0].client),
	}

	tests := []struct {
		required []string
		want     Weed
	}{
		{required: []string{"ec2:DescribeInstances", "iam:CreateAccessKey"}, want: Weed{}},
		{required: []string{}, want: Weed{Unnecessary: []string{"ec2:DescribeInstances", "iam:CreateAccessKey"}}},
		{required: []string{"ec2:DescribeInstances"}, want: Weed{Unnecessary: []string{"iam:CreateAccessKey"}}},
		{required: []string{"ec2:DescribeInstances", "iam:CreateAccessKey", "ec2:AssociateIamInstanceProfile"}, want: Weed{Missing: []string{"ec2:AssociateIamInstanceProfile"}}},
		{required: []string{"ec2:DescribeInstances", "ec2:AssociateIamInstanceProfile"}, want: Weed{Missing: []string{"ec2:AssociateIamInstanceProfile"}, Unnecessary: []string{"iam:CreateAccessKey"}}},
	}

	for _, test := range tests {
		got, err := weedClient.Find(test.required, iamRole)
		if err != nil {
			t.Errorf("error finding weed: %v", err)
		}

		if cmp.Equal(got, test.want) == false {
			t.Errorf("weedClient.Find(%v, %v) = %v, want %v", test.required, iamRole, got, test.want)
		}
	}
}
