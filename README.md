# Weed - IAM role diff checker

## Installation
### Required IAM Policy
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ListRolePermissions",
      "Action": [
        "iam:GetPolicy",
        "iam:GetPolicyVersion",
        "iam:GetRolePolicy",
        "iam:ListAttachedRolePolicies",
        "iam:ListRolePolicies"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```

## Todo
- [*] Diff checker
- [*] CLI Wrapper
- [*] IAM enforcer
- [*] GO Release
- [ ] README / Docs
- [ ] Examples
- [ ] Handle * permissions
