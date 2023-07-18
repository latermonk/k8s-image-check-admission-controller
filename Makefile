.PHONY: pull all docker



all: pull docker

pull:
	git clone https://github.com/latermonk/k8s-image-check-admission-controller.git
#    cd k8s-image-check-admission-controller

docker:
	docker build . -t ibackchina2018/k8s-image-admission-controller:latest