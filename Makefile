SHELL := /bin/bash

ROOT := $(shell git rev-parse --show-toplevel)

VERSION ?= $(shell git describe --dirty="-dev")

DOCKER_IMG ?= form3tech/cilium-etcd-watchdog
DOCKER_TAG ?= $(VERSION)

.PHONY: docker.build
docker.build:
	docker build -t $(DOCKER_IMG):$(DOCKER_TAG) $(ROOT)

.PHONY: docker.push
docker.push: docker.build
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
	docker push $(DOCKER_IMG):$(DOCKER_TAG)
