# Weed - IAM role diff checker

## Installation
### 1. Admission Controller Certificate Creation
```bash
brew install openssl@1.1 # Needed for creating certs with SANs
mkdir certs

export ENFORCER_SERVICE_NAME=iam-enforcer
export ENFORCER_TLS_SECRET_NAME=iam-enforcer-tls # Set it in the helm chart value file for `tlsCertSecretName`
export ENFORCER_K8S_NAMESPACE=default # Should be the namespace you'll install the helm chart to

# Creating certificates
openssl genrsa -out certs/ca.key 2048
openssl req -new -x509 -sha256 -key certs/ca.key -days 365 -subj "/CN=IAM Enforcer Admission Controller" -out certs/ca.crt
openssl genrsa -out certs/admission-tls.key 2048
openssl req -new -key certs/admission-tls.key -subj "/CN=${ENFORCER_SERVICE_NAME}.${ENFORCER_K8S_NAMESPACE}.svc"  | openssl x509 -req -sha256 -days 365 \
  -extfile <(printf "subjectAltName=DNS:${ENFORCER_SERVICE_NAME}.${ENFORCER_K8S_NAMESPACE}.svc,DNS:${ENFORCER_SERVICE_NAME}.${ENFORCER_K8S_NAMESPACE}.svc") \
  -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial \
  -out certs/admission-tls.crt

# Creating k8s Secret with TLS certs
kubectl create secret tls -n $ENFORCER_K8S_NAMESPACE $ENFORCER_TLS_SECRET_NAME \
    --cert "certs/admission-tls.crt" \
    --key "certs/admission-tls.key"

cat certs/ca.crt | base64 | tr -d '\n' > certs/ca_bundle.txt # Set it in the helm chart value file for `caBundle`
```

### 2. IAM Policy
Create an IAM role with the following policy:
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

### 3. Deploy Helm Chart
```bash
export VALUES_FILE=<values-file>
echo "tlsCertSecretName: $ENFORCER_TLS_SECRET_NAME" > $VALUES_FILE
echo "caBundle: $(`cat certs/ca_bundle.txt`)" >> $VALUES_FILE

# Assign created IAM role:
# On EKS use ServiceAccount:
# echo "serviceAccount:" >> $VALUES_FILE
# echo "  name: <sa-name>" >> $VALUES_FILE
#
# Using kube2iam use pod spec annotation:
# echo "podAnnotations:" >> $VALUES_FILE
# echo "  iam.amazonaws.com/role: <aws-role-arn>" >> $VALUES_FILE

# Override AWS region:
# echo "environmentVariables:" >> $VALUES_FILE
# echo "  - name: \"AWS_REGION\"" >> $VALUES_FILE
# echo "    value: \"<region>\"" >> $VALUES_FILE

helm install iam-enforcer ./helm --values $VALUES_FILE
```

## Todo
- [x] Diff checker
- [x] CLI Wrapper
- [x] IAM enforcer
- [x] GO Release
- [x] README / Docs
- [ ] Examples
- [ ] Handle * permissions
