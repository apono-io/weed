package weed

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type IIamClient interface {
	ListRolePolicies(ctx context.Context, params *iam.ListRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListRolePoliciesOutput, error)
	GetRolePolicy(ctx context.Context, params *iam.GetRolePolicyInput, optFns ...func(*iam.Options)) (*iam.GetRolePolicyOutput, error)
	ListAttachedRolePolicies(ctx context.Context, params *iam.ListAttachedRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error)
	GetPolicy(ctx context.Context, params *iam.GetPolicyInput, optFns ...func(*iam.Options)) (*iam.GetPolicyOutput, error)
	GetPolicyVersion(ctx context.Context, params *iam.GetPolicyVersionInput, optFns ...func(*iam.Options)) (*iam.GetPolicyVersionOutput, error)
}

type AwsService struct {
	iamClient IIamClient
	ctx       context.Context
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

type AWSPolicyPermissionStatement struct {
	Effect   string         `json:"Effect"`
	Action   ArrayOfStrings `json:"Action"`
	Resource ArrayOfStrings `json:"Resource"`
	// Condition interface{}    `json:"Condition"` // TODO: add logic to handle this
}

type PolicyPermission struct {
	Version   string                         `json:"Version"`
	Statement []AWSPolicyPermissionStatement `json:"Statement"`
}

func NewAwsService() (service *AwsService, err error) {
	service.ctx = context.TODO()

	cfg, err := config.LoadDefaultConfig(service.ctx)
	if err != nil {
		return service, fmt.Errorf("error loading config: %v", err)
	}

	// Create a IAM client from config
	service.iamClient = iam.NewFromConfig(cfg)

	return
}

func (svc AwsService) RolePermissions(roleName string) (permissions []string, err error) {
	rolePolicies, err := svc.iamClient.ListRolePolicies(svc.ctx, &iam.ListRolePoliciesInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return permissions, fmt.Errorf("error listing role policies: %v", err)
	}

	for _, pn := range rolePolicies.PolicyNames {
		policy, err := svc.iamClient.GetRolePolicy(svc.ctx, &iam.GetRolePolicyInput{
			PolicyName: aws.String(pn),
			RoleName:   aws.String(roleName),
		})
		if err != nil {
			return permissions, fmt.Errorf("error getting role policy: %v", err)
		}

		policyPermissions, err := svc.permissionsFromDocument(policy.PolicyDocument)
		if err != nil {
			return permissions, fmt.Errorf("error getting role policy statements: %v", err)
		}

		permissions = append(permissions, policyPermissions...)
	}

	attachedRolePolicies, err := svc.iamClient.ListAttachedRolePolicies(svc.ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return permissions, fmt.Errorf("error listing attached role policies: %v", err)
	}

	for _, ap := range attachedRolePolicies.AttachedPolicies {
		policy, err := svc.iamClient.GetPolicy(svc.ctx, &iam.GetPolicyInput{
			PolicyArn: ap.PolicyArn,
		})
		if err != nil {
			return permissions, fmt.Errorf("error getting attached role policy: %v", err)
		}

		pv, err := svc.iamClient.GetPolicyVersion(svc.ctx, &iam.GetPolicyVersionInput{
			PolicyArn: policy.Policy.Arn,
			VersionId: policy.Policy.DefaultVersionId,
		})
		if err != nil {
			return permissions, fmt.Errorf("error getting attached role policy version: %v", err)
		}

		policyPermissions, err := svc.permissionsFromDocument(pv.PolicyVersion.Document)
		if err != nil {
			return permissions, fmt.Errorf("error getting role policy statements: %v", err)
		}

		permissions = append(permissions, policyPermissions...)
	}

	return
}

func (svc AwsService) permissionsFromDocument(document *string) (permissions []string, err error) {
	if document == nil {
		return permissions, nil
	}

	decodedValue, err := url.QueryUnescape(*document)
	if err != nil {
		return permissions, fmt.Errorf("error decoding policy document: %v", err)
	}

	var policyPermission PolicyPermission
	if err := json.Unmarshal([]byte(decodedValue), &policyPermission); err != nil {
		return permissions, fmt.Errorf("error unmarshaling policy document: %v", err)
	}

	var statements []AWSPolicyPermissionStatement
	for _, statement := range policyPermission.Statement {
		if statement.Effect == "Allow" {
			statements = append(statements, statement)
		}
	}

	permissionsMap := make(map[string]interface{})

	for _, s := range statements {
		for _, ps := range s.Action {
			permissionsMap[ps] = true
		}
	}

	for s := range permissionsMap {
		permissions = append(permissions, s)
	}

	return
}
