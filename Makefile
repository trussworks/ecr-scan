SHELL = /bin/sh
VERSION = 0.0.1

.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: local_build
local_build: ## Build ecr-scan locally
	@./scripts/local-build

.PHONY: lambda_build
lambda_build: ## Build lambda binary
	@./scripts/lambda-build $(VERSION)

.PHONY: lambda_release
lambda_release: lambda_build ## Release lambda zip file to S3
	@./scripts/lambda-release $(S3_BUCKET) $(VERSION)

.PHONY: lambda_run
lambda_run: lambda_build ## Run the lambda handler in docker
	@./scripts/lambda-run

.PHONY: clean
clean: ## Clean all generated files
	rm -rf ./bin
	rm -rf ./dist

default: help
