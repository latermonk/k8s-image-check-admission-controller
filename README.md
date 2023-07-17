# k8s-image-check-admission-controller
k8s-image-check-admission-controller

# BUILD
```shell
git clone https://github.com/latermonk/k8s-image-check-admission-controller.git
cd k8s-image-check-admission-controller
```

```shell
rm -rf go.*
go mod init k8s-image-check-admission-controller
go mod tidy && go run main.go
```


```shell
docker build . -t ibackchina2018/k8s-image-admission-controller:latest
docker login
docker push ibackchina2018/k8s-image-admission-controller:latest
```


---
# Test

## Cert-Manager
```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml

```
## cmdtool:
```shell
OS=$(go env GOOS); ARCH=$(go env GOARCH); curl -fsSL -o cmctl.tar.gz https://github.com/cert-manager/cert-manager/releases/latest/download/cmctl-$OS-$ARCH.tar.gz
tar xzf cmctl.tar.gz
sudo install cmctl /usr/local/bin
cmctl check api
```


