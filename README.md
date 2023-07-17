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

##  Install resources
```bash
kubectl apply -f k8s/00_namespace.yaml
kubectl apply -f k8s/10_ca_certificate.yaml
kubectl apply -f k8s/10_certificate.yaml
kubectl apply -f k8s/20_deployment.yaml
kubectl apply -f k8s/20_service.yaml
kubectl apply -f k8s/30_validatingwebhookconfiguration.yaml
```

##  Test that pod is denied
```bash
kubectl apply -f k8s/90_pod-test.yaml
```
