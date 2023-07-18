.PHONY: all build docker cert test



all:  build docker cert


build:
	rm -rf  go.* && \
	go mod init k8s-image-check-admission-controller && \
	go mod tidy && \
	go build -o  abc   ./cmd/k8s-image-admission-controller/main.go

docker:
	docker build . -t ibackchina2018/k8s-image-admission-controller:latest


cert:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml && \
	OS=$(go env GOOS); ARCH=$(go env GOARCH); curl -fsSL -o cmctl.tar.gz https://github.com/cert-manager/cert-manager/releases/latest/download/cmctl-$OS-$ARCH.tar.gz && \
    tar xzf cmctl.tar.gz && \
    sudo install cmctl /usr/local/bin  && \
    cmctl check api

test:
	kubectl apply -f 01-backend/00_namespace.yaml && \
    kubectl apply -f 01-backend/10_ca_certificate.yaml && \
    kubectl apply -f 01-backend/10_certificate.yaml && \
    kubectl apply -f 01-backend/20_deployment.yaml && \
    kubectl apply -f 01-backend/20_service.yaml