# Kind: https://kind.sigs.k8s.io run K8s locally
# Kubectl: https://kubernetes.io/docs/reference/kubectl/overview/ control K8s

SHELL := /bin/bash

# Run the app as is.
run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

# Generate the private and public pem files.
admin:
	go run app/tooling/admin/main.go


# ============================================================================================================
# Running Tests on local

# ./... = step through whole project
# -count=1 = ignore cache and run tests every time
# staticcheck = make sure we consider linting in our tests
test:
	go test ./... -v -count=1
	staticcheck -checks=all ./...

# ============================================================================================================
# Testing running system

# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem

# Brings up a GUI in your terminal for the metrics of the current service.
# Make sure to run metrics from the root of the project or copy and paste the command into terminal.
# You will need the following: https://github.com/divan/expvarmon
metrics: 
	expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# For testing load on the service.
# You will need the following: https://github.com/rakyll/hey
load-test:
	hey -m GET -c 100 -n 10000 http://localhost:3000/v1/test

# ============================================================================================================
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

# ============================================================================================================
# Running from within k8s/kind

KIND_CLUSTER := ryan-starter-cluster

# Build the image, start the cluster and then load the image into the cluster.
kind-start: all kind-up kind-load kind-apply

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
# we will apply the database before the sales service
kind-apply:
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -

kind-restart: 
	kubectl -n sales-system rollout restart deployment sales-pod

# Only need to run this when change application logic
kind-update: all kind-load kind-restart

# load in the new image to kind and then apply it. Use when make changes to k8s, Docker.
kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=sales

# get the status of the pods
kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-sales:
	kubectl get pods -o wide --watch 
kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

# get the logs for the service
kind-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100 --namespace=sales-system


# ============================================================================================================
# Module Support

tidy:
	go mod tidy
	go mod vendor