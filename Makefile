.PHONY: pull all docker



all: pull docker buildbianry

pull:
	git clone https://github.com/latermonk/k8s-image-check-admission-controller.git \
    cd k8s-image-check-admission-controller

build:
	rm -rf  go.* \
	go mod init k8s-image-check-admission-controller \
	go mod tidy && go build -o  abc   ./cmd/k8s-image-admission-controller/main.go

docker:
	docker build . -t ibackchina2018/k8s-image-admission-controller:latest