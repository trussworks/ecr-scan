SHELL = /bin/sh
VERSION = 0.0.1

.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build_local
build_local: ## Build ecr-scan locally
	go build -o bin/ecr-scan .

.PHONY: build_lambda_builder
build_lambda_builder: ## Build docker image for building lambda handler
	docker image pull lambci/lambda:build-go1.x
	docker image build -t ecr-scan-builder:"${VERSION}" .

.PHONY: build_lambda_handler
build_lambda_handler: build_lambda_builder ## Build lambda binary
	docker container run --rm -it -v "${PWD}":/app ecr-scan-builder:"${VERSION}" go build -o bin/ecr-scan .

.PHONY: run_lambda_handler
run_lambda_handler: build_lambda_handler ## Run the lambda handler in the background
	docker container run --rm -e LAMBDA=1 -e DOCKER_LAMBDA_STAY_OPEN=1 -p 9001:9001 -v "${PWD}":/var/task:ro,delegated lambci/lambda:go1.x bin/ecr-scan

.PHONY: clean
clean: ## Clean all generated files
	rm -rf ./bin
	rm -rf ./dist

default: help
