# Kind: https://kind.sigs.k8s.io run K8s locally
# Kubectl: https://kubernetes.io/docs/reference/kubectl/overview/ control K8s

SHELL := /bin/bash

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

# ======================================================
# Build Containers

VERSION := 1.0

all: sales-api

# Build the docker image.
sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:$(VERSION) \
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

# navigate into the kind/sales-pod and edit the image name to include the version
# load our local images into the kind environment
kind-load:
	cd zarf/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

# Tell K8s to apply the namespace to the deployment
# kustomize will build the final yaml starting from the service-pod
kind-apply:
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -

kind-update: all kind-load kind-restart

kind-restart: 
	kubectl rollout restart deployment sales-pod

# load in the new image to kind and then apply it
kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=sales

# get the status of the pods
kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

# get the logs for the service
kind-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100 --namespace=sales-system | go run app/tooling/logfmt/main.go

# ======================================================
# Module Support

tidy:
	go mod tidy
	go mod vendor