SHELL = /bin/sh

help:  ## Print the help documentation
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

check_git_status: ## Check for modified, added, and unstaged files
	@./scripts/check-git-status

test: ## Run unit tests
	@./scripts/make-test

local_build: test clean ## Build ecr-scan locally
	@./scripts/local-build

lambda_build: test clean ## Build lambda binary
	@./scripts/lambda-build

lambda_release: check_git_status lambda_build ## Release lambda zip file to S3
	@./scripts/lambda-release $(S3_BUCKET)

lambda_run: lambda_build ## Run the lambda handler in docker
	@./scripts/lambda-run

clean: ## Clean all generated files
	rm -rf ./bin
	rm -rf ./dist

default: help

.PHONY: help check_git_status test local_build lambda_build lambda_release lambda_run clean
