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
go mod tidy && go build -o  abc   ./cmd/k8s-image-admission-controller/main.go
```


```shell
docker build . -t ibackchina2018/k8s-image-admission-controller:latest

```

```shell
docker login

```

```shell
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

```shell
git clone https://github.com/latermonk/k8s-image-check-admission-controller.git
cd k8s-image-check-admission-controller
```

## Backend svc
```bash
kubectl apply -f 01-backend/00_namespace.yaml
kubectl apply -f 01-backend/10_ca_certificate.yaml
kubectl apply -f 01-backend/10_certificate.yaml
kubectl apply -f 01-backend/20_deployment.yaml
kubectl apply -f 01-backend/20_service.yaml

```
## webhook
```shell
kubectl apply -f 02-validatingwebhookconfiguration.yaml
```


##  Test that pod is denied
```bash
kubectl apply -f 03_pod-test.yaml
```



---
#  Install golang 1.20

```shell
wget https://go.dev/dl/go1.20.6.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.6.linux-amd64.tar.gz

```

```shell
export PATH=$PATH:/usr/local/go/bin
```

# test go version
```shell
go version
```

---

#  Dive

```shell
wget https://github.com/wagoodman/dive/releases/download/v0.11.0/dive_0.11.0_linux_amd64.tar.gz && tar -zxvf dive_0.11.0_linux_amd64.tar.gz && install dive /usr/local/bin 
```

Test:    
```shell
CI=true dive ibackchina2018/ubuntu-sshd:5g
```