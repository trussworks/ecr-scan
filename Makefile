SHELL = /bin/sh
VERSION = 0.0.1

help:  ## Print the help documentation
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test:
	@./scripts/make-test

local_build: test ## Build ecr-scan locally
	@./scripts/local-build

lambda_build: test ## Build lambda binary
	@./scripts/lambda-build $(VERSION)

lambda_release: lambda_build ## Release lambda zip file to S3
	@./scripts/lambda-release $(S3_BUCKET) $(VERSION)

lambda_run: lambda_build ## Run the lambda handler in docker
	@./scripts/lambda-run

clean: ## Clean all generated files
	rm -rf ./bin
	rm -rf ./dist

default: help

.PHONY: help test local_build lambda_build lambda_release lambda_run clean
