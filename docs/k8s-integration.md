# Kubernetes Integration
Below are the installation instructions for IAM Enforcer admission controller.

## Usage
IAM Enforcer uses pod spec annotations to check the required AWS IAM actions that should be validation.
To get the name of the role that is assigned to the pod we look for one of the following annotations in this order:
1. `iam-enforcer.apono.io/role-arn`
2. `iam.amazonaws.com/role`
3. `eks.amazonaws.com/role-arn` (this annotation is looked for in the service account)

To get the expected AWS IAM actions iam-enforcer looks for one of the following annotations in this order:
1. `iam-enforcer.apono.io/required-actions` - Comma separated list of IAM actions like: `ec2:DescribeInstances,ec2:DescribeVpcs`

Example:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: default
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
      annotations:
        iam.amazonaws.com/role: arn:aws:iam::123456789012:role/my-role
        iam-enforcer.apono.io/required-actions: ec2:DescribeInstances,ec2:DescribeVpcs
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
```

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
