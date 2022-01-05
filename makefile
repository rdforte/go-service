# Kind: https://kind.sigs.k8s.io run K8s locally
# Kubectl: https://kubernetes.io/docs/reference/kubectl/overview/ control K8s

SHELL := /bin/bash

tidy:
	go mod tidy
	go mod vendor

run:
	go run main.go

# ======================================================
# Build Containers

VERSION := 1.0

all: service

service:
	docker build \
		-f zarf/docker/dockerfile \
		-t service-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ======================================================
# Running from within k8s/kind

KIND_CLUSTER := ryan-starter-cluster

# Kind release used for our project: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.1
# The image used below was copied by the above link and supports both amd64 and arm64.
kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

# load our local images into the kind environment
kind-load:
	kind load docker-image service-amd64:$(VERSION) --name $(KIND_CLUSTER)

# Tell K8s to apply the namespace to the deployment
kind-apply:
	cat zarf/k8s/base/service-pod/base-service.yaml | kubectl apply -f -

kind-update: all kind-load kind-restart

# get the status of the pods
kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

# get the logs for the service
kind-logs:
	kubectl logs -l app=service --all-containers=true -f --tail=100 --namespace=service-system