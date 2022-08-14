# WEED - CLI Tool
Below are the installation instructions for WEED CLI tool.

## Usage
WEED can be used to compare policy file or a list of actions to IAM role in AWS.
AWS role ARN/name must be passed to the command using `--role` flag.
Policy file comparison or action list can be used with `--policy-file` or `--actions` flags.
The tool can exit with status code 1 on missing actions by passing `--fail-on-missing` flag or on any difference using `--fail-on-diff` flag.

### Example
Compare IAM role to a policy file:
```shell
weed --role-arn weed-demo --policy-file policy.json --fail-on-missing
```

Compare IAM role to a list of actions:
```shell
weed --role-arn weed-demo --actions ec2:DescribeInstances,ec2:DescribeTags --fail-on-diff
```

## Installation
1. Go to the [releases page](https://github.com/apono-io/weed/releases/latest) and download the archive for your environment.
2. Extract the weed binary from the archive and move it to a directory within your environment's PATH.
3. Run `weed version` to see the installed version.
