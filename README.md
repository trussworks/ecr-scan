# DEPRECATION NOTICE

This repository has been deprecated and is no longer maintained. Should you need
to continue to use it, please fork the repository. Thank you.

As alternatives, we recommend that project teams consider Trivy and/or Grype for
image vulnerability scans.

## ecr-scan

## Description

`ecr-scan` is used to scan and retrieve findings for images in ECR.

Run it from the command-line or as a Lambda function.

## Installation

```shell
go get -u github.com/trussworks/ecr-scan
```

## Usage

```shell
ecr-scan is an application for analyzing ECR scan findings

Usage:
  ecr-scan [flags]

Flags:
  -h, --help                help for ecr-scan
      --lambda              Run as Lambda function
  -m, --maxScanAge int      Maximum allowed age for image scan (hours) (default 24)
      --profile string      The AWS profile to use
      --region string       The AWS region to use
  -r, --repository string   ECR repository where the image is located
  -t, --tag string          Image tag to retrieve findings for
```

If no scan is found, or if the most recent scan exceeds `maxScanAge`, `ecr-scan`
will re-scan the image.

## Examples

Run the command like this:

```shell
ecr-scan --profile app-dev --region us-west-2 --repository app-ecr --tag 50e9216704a67c97664dbbac521b3a674c61cee9
```

## Developer Setup

### Available options

```shell
$ make help
check_git_status               Check for modified, added, and unstaged files
clean                          Clean all generated files
help                           Print the help documentation
lambda_build                   Build lambda binary
lambda_release                 Release lambda zip file to S3
lambda_run                     Run the lambda handler in docker
local_build                    Build ecr-scan locally
pre_commit                     Run all pre-commit checks
test                           Run unit tests
```

### Release deployment package

```shell
make S3_BUCKET=your-s3-bucket lambda_release
```
