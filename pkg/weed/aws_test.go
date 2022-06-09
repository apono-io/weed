package weed

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/google/go-cmp/cmp"
)

type mockIamClient struct {
	iam.Client
	ListRolePoliciesOutput         iam.ListRolePoliciesOutput
	GetRolePolicyOutput            iam.GetRolePolicyOutput
	ListAttachedRolePoliciesOutput iam.ListAttachedRolePoliciesOutput
	GetPolicyOutput                iam.GetPolicyOutput
	GetPolicyVersionOutput         iam.GetPolicyVersionOutput
}

func (m mockIamClient) ListRolePolicies(ctx context.Context, params *iam.ListRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListRolePoliciesOutput, error) {
	return &m.ListRolePoliciesOutput, nil
}

func (m mockIamClient) GetRolePolicy(ctx context.Context, params *iam.GetRolePolicyInput, optFns ...func(*iam.Options)) (*iam.GetRolePolicyOutput, error) {
	return &m.GetRolePolicyOutput, nil
}

func (m mockIamClient) ListAttachedRolePolicies(ctx context.Context, params *iam.ListAttachedRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error) {
	return &m.ListAttachedRolePoliciesOutput, nil
}

func (m mockIamClient) GetPolicy(ctx context.Context, params *iam.GetPolicyInput, optFns ...func(*iam.Options)) (*iam.GetPolicyOutput, error) {
	return &m.GetPolicyOutput, nil
}

func (m mockIamClient) GetPolicyVersion(ctx context.Context, params *iam.GetPolicyVersionInput, optFns ...func(*iam.Options)) (*iam.GetPolicyVersionOutput, error) {
	return &m.GetPolicyVersionOutput, nil
}

var awsTestCases = []struct {
	client   mockIamClient
	expected []string
}{
	{
		client: mockIamClient{
			ListRolePoliciesOutput: iam.ListRolePoliciesOutput{
				PolicyNames: []string{"test-policy-1"},
			},
			GetRolePolicyOutput: iam.GetRolePolicyOutput{
				PolicyDocument: aws.String("%7B%20%22Version%22%3A%20%222012-10-17%22%2C%20%22Statement%22%3A%20%5B%7B%20%22Effect%22%3A%20%22Allow%22%2C%20%22Action%22%3A%20%5B%22ec2%3ADescribeInstances%22%5D%2C%20%22Resource%22%3A%20%22*%22%20%7D%5D%20%7D"),
			},
			ListAttachedRolePoliciesOutput: iam.ListAttachedRolePoliciesOutput{
				AttachedPolicies: []types.AttachedPolicy{
					{
						PolicyArn:  aws.String("arn:aws:iam::aws:policy/test-policy-2"),
						PolicyName: aws.String("test-policy-2"),
					},
				},
			},
			GetPolicyOutput: iam.GetPolicyOutput{
				Policy: &types.Policy{
					Arn:              aws.String("arn:aws:iam::aws:policy/test-policy-1"),
					DefaultVersionId: aws.String("v1"),
				},
			},
			GetPolicyVersionOutput: iam.GetPolicyVersionOutput{
				PolicyVersion: &types.PolicyVersion{
					Document: aws.String("%7B%20%22Version%22%3A%20%222012-10-17%22%2C%20%22Statement%22%3A%20%5B%7B%20%22Effect%22%3A%20%22Allow%22%2C%20%22Action%22%3A%20%5B%22iam%3ACreateAccessKey%22%5D%2C%20%22Resource%22%3A%20%22*%22%20%7D%5D%20%7D"),
				},
			},
		},
		expected: []string{"ec2:DescribeInstances", "iam:CreateAccessKey"},
	},
}

func newMockAwsService(client mockIamClient) *AwsService {
	return &AwsService{
		ctx:       context.Background(),
		iamClient: client,
	}
}

func TestRolePermissions(t *testing.T) {
	for _, c := range awsTestCases {
		service := newMockAwsService(c.client)

		got, err := service.RolePermissions("arn:aws:iam::123456789012:role/test-role")
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if cmp.Equal(got, c.expected) == false {
			t.Errorf("got %v, want %v", got, c.expected)
		}
	}
}
