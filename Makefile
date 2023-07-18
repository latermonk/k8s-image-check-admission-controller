.PHONY: pull all docker



all: pull docker

pull:
	git clone https://github.com/latermonk/k8s-image-check-admission-controller.git
#    cd k8s-image-check-admission-controller

docker:
	docker build --pull --build-arg COSIGN_VERSION=$(COSIGN_VERSION) -f docker/Dockerfile -t $(IMAGE_REPOSITORY):v$(VERSION) .