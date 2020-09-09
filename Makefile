
SHELL := $(shell which bash)
OSTYPE := $(shell uname)
DOCKER := $(shell command -v docker)
GID := $(shell id -g)
UID := $(shell id -u)
VERSION ?= $(shell git describe --tags --always)

UNIT_TEST_CMD := ./scripts/check/unit-test.sh
INTEGRATION_TEST_CMD := ./scripts/check/integration-test.sh
CHECK_CMD := ./scripts/check/check.sh

DEV_IMAGE_NAME := pelotech/templated-configuration-webhook-dev
WEBHOOK_IMAGE_NAME ?=  pelotech/templated-configuration-webhook
ENVTEMPLATE_IMAGE_NAME ?= pelotech/envtemplate_webhook

DOCKER_RUN_CMD := docker run --env ostype=$(OSTYPE) -v ${PWD}:/src --rm -it ${DEV_IMAGE_NAME}
BUILD_BINARY_CMD := VERSION=${VERSION} ./scripts/build/build.sh
BUILD_DEV_IMAGE_CMD := IMAGE=${DEV_IMAGE_NAME} DOCKER_FILE_PATH=./docker/dev/Dockerfile VERSION=latest ./scripts/build/build-image.sh
BUILD_WEBHOOK_IMAGE_CMD := IMAGE=${PROD_IMAGE_NAME} DOCKER_FILE_PATH=./docker/prod/Dockerfile VERSION=${VERSION} ./scripts/build/build-image.sh
BUILD_ENVTEMPLATE_IMAGE_CMD := IMAGE=${ENVTEMPLATE_IMAGE_NAME} DOCKER_FILE_PATH=./docker/envtemplate/Dockerfile VERSION=${VERSION} ./scripts/build/build-image.sh
PUBLISH_WEBHOOK_IMAGE_CMD := IMAGE=${WEBHOOK_IMAGE_NAME} VERSION=${VERSION} ./scripts/build/publish-image.sh
PUBLISH_ENVTEMPLATE_IMAGE_CMD := IMAGE=${ENVTEMPLATE_IMAGE_NAME} VERSION=${VERSION} ./scripts/build/publish-image.sh
GEN_CERTS_CMD := ./scripts/gen-certs.sh


help: ## Show this help
	@echo "Help"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[93m %s\n", $$1, $$2}'

.PHONY: default
default: help

.PHONY: build-webhook-image
build-webhook-image: ## Builds the webhook docker image.
	@$(BUILD_WEBHOOK_IMAGE_CMD)
	
.PHONY: build-envtemplate-image
build-envtemplate-image: ## Builds the envtemplate docker image.
	@$(BUILD_ENVTEMPLATE_IMAGE_CMD)

.PHONY: publish-webhook-image
publish-webhook-image: ##Publishes the webhook docker image.
	@$(PUBLISH_WEBHOOK_IMAGE_CMD)

.PHONY: publish-envtemplate-image
publish-envtemplate-image: ##Publishes the envtemplate docker image.
	@$(PUBLISH_ENVTEMPLATE_IMAGE_CMD)

.PHONY: build-dev-image
build-dev-image:  ## Builds the development docker image.
	@$(BUILD_DEV_IMAGE_CMD)

build: build-dev-image ## Builds the production binary.
	@$(DOCKER_RUN_CMD) /bin/sh -c '$(BUILD_BINARY_CMD)'

.PHONY: test
test: build-dev-image  ## Runs unit test.
	@$(DOCKER_RUN_CMD) /bin/sh -c '$(UNIT_TEST_CMD)'

.PHONY: check
check: build-dev-image  ## Runs checks.
	@$(DOCKER_RUN_CMD) /bin/sh -c '$(CHECK_CMD)'

.PHONY: integration
integration: build-dev-image ## Runs integration test.
	@$(DOCKER_RUN_CMD) /bin/sh -c '$(INTEGRATION_TEST_CMD)'

.PHONY: deps
deps:  ## Fixes the dependencies
	@$(DOCKER_RUN_CMD) /bin/sh -c './scripts/deps.sh'

.PHONY: ci-unit-test
ci-test:  ## Runs unit test in CI environment (without docker).
	@$(UNIT_TEST_CMD)

.PHONY: ci-check
ci-check:  ## Runs checks in CI environment (without docker).
	@$(CHECK_CMD)

.PHONY: ci-integration-test
ci-integration: ## Runs integraton test in CI environment (without docker).
	@$(INTEGRATION_TEST_CMD)

.PHONY: gen-deploy-certs
gen-deploy-certs: ## Generate deploy files webhook certificates.
	@$(GEN_CERTS_CMD)
