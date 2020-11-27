VERSION ?= 0.2.0
IMG=agill17/pod-mutating-webhook:${VERSION}

build: docker-build docker-push

docker-build:
	go fmt .
	go vet .
	docker build . -t ${IMG}

docker-push:
	docker push ${IMG}

install:
	./setup.sh pod-mutating-webhook default
	helm upgrade -i pod-mutating-webhook pod-mutating-webhook -f pod-mutating-webhook/tls.values.yaml

uninstall:
	helm uninstall pod-mutating-webhook
