SHELL = /bin/sh
VERSION = 0.0.1

.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: local_build
local_build: ## Build ecr-scan locally
	go build -o bin/ecr-scan .

.PHONY: lambda_build
lambda_build: ## Build lambda binary
	docker image build -t ecr-scan-builder:"$(VERSION)" .
	docker container run --rm -it -v "$(PWD)":/app ecr-scan-builder:$(VERSION) go build -o bin/ecr-scan .

.PHONY: lambda_release
lambda_release: lambda_build ## Release lambda zip file to S3
	zip -j ecr-scan.zip ./bin/ecr-scan
	#aws s3 cp --sse AES256 ecr-scan.zip s3://"$(S3_BUCKET)"/ecr-scan/"$(VERSION)"/

.PHONY: run_lambda
run_lambda: lambda_build ## Run the lambda handler in docker
	docker container run --rm -e LAMBDA=1 -e DOCKER_LAMBDA_STAY_OPEN=1 -p 9001:9001 -v "$(PWD)":/var/task:ro,delegated lambci/lambda:go1.x bin/ecr-scan

.PHONY: clean
clean: ## Clean all generated files
	rm -rf ./bin
	rm -rf ./dist

default: help
