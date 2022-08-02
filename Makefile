IMAGE_NAME ?= quay.io/fhke97/kubectl-draincheck
IMAGE_TAG  ?= $(shell cat ./VERSION)

IMAGE = $(IMAGE_NAME):$(IMAGE_TAG)
IMAGE_LATEST = $(IMAGE_NAME):latest

.DEFAULT_GOAL = help

.PHONY: help
help: ## Display help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build docker image
	docker build -t $(IMAGE) -t $(IMAGE_LATEST) .

.PHONY: publish
publish: build ## Build docker image & push to remote repo
	docker push $(IMAGE)
	docker push $(IMAGE_LATEST)
