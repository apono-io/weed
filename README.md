# WEED - What EEs Different - IAM role diff checker
[![Maintained by Apono.io](https://img.shields.io/badge/maintained%20by-apono.io-3f9fcc)](https://apono.io/?utm_source=github&utm_medium=organic_oss&utm_campaign=weed)
[![Build Status](https://github.com/apono-io/weed/workflows/go/badge.svg)](https://github.com/apono-io/weed/actions?query=workflow%3Ago)
![Go version](https://img.shields.io/github/go-mod/go-version/apono-io/weed)
![GitHub Release (latest)](https://img.shields.io/github/v/release/apono-io/weed)

## Prevent runtime errors in production ahead of time!
![Introduction](.github/assets/intro.jpg)

Have you ever pushed to production only to find out that the permissions between staging and production environments are out of sync? resulting in access errors in a live environment? 
Well we have, we got frustrated, learned and created WEED!!!

## What is WEED?

WEED is a CLI tool that assures permissions are synced between different environments.
WEED Checks for permission differences between requested permissions in an environment to current environment.
WEED maps permissions on both environments checking for discrepancies that might cause access errors in production

### Components

***WEED CLI***  -
Discovers Diff in permissions between environments to avoid those pesky 403 errors in production.
Can be used to verify permissions manually or as part of the CI CD cycle.

***IAM Enforcer – Kubernetes Admission Controller*** –
Intercepts API requests to k8s api-server and acts as a validation layer, assuring all requested permissions are available before applying the changes.

### Prerequisites

- AWS Account
- Role in AWS

## Getting Started
Installation instruction for the Kubernetes integrations are available [here](docs/k8s-integration.md), instruction for the CLI tool are available [here](docs/cli.md).

![Demo](.github/assets/demo.gif)
