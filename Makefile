IMAGE_NAME ?= quay.io/fhke97/kubectl-draincheck
VERSION  ?= $(shell cat ./VERSION)

IMAGE = $(IMAGE_NAME):$(VERSION)
IMAGE_LATEST = $(IMAGE_NAME):latest

GO_VERSION ?= 1.18
GO_IMAGE ?= golang

.DEFAULT_GOAL = help

.PHONY: help
help: ## Display help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build docker image
	docker build --build-arg 'GO_IMAGE=$(GO_IMAGE)' --build-arg 'GO_VERSION=$(GO_VERSION)' -t $(IMAGE) -t $(IMAGE_LATEST) .

.PHONY: publish
publish: build ## Build docker image & push to remote repo
	docker push $(IMAGE)
	docker push $(IMAGE_LATEST)

.PHONY: tag_release
tag_release: ## Create a git tag for the current release
	git tag $(VERSION)


.PHONY: test
test: ## Run tests
	go test -v ./...