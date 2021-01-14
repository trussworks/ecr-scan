# ecr-scan

## Description

`ecr-scan` is used to scan and retrieve findings for images in ECR.

Run it from the command-line or as a Lambda function.

## Installation

```sh
go get -u github.com/trussworks/ecr-scan
```

## Usage

```sh
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

```sh
ecr-scan --profile app-dev --region us-west-2 --repository app-ecr --tag 50e9216704a67c97664dbbac521b3a674c61cee9
```
