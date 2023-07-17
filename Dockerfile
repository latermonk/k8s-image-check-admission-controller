ARG APP_NAME=k8s-image-admission-controller
ARG USER=appuser

############################
# STEP 1 build executable binary
############################
FROM golang:1.20-alpine as builder
ENV UID=10001
ARG USER
ARG APP_NAME

# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
# hadolint ignore=DL3018
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
RUN go mod download && go mod verify && GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/main cmd/${APP_NAME}/main.go

############################
# STEP 2 build a small image
############################
FROM scratch
ARG USER
ARG APP_NAME
LABEL org.opencontainers.image.description="Secure image for docker Golang bootstrap app"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/bin/main /app/bin/main

EXPOSE 8080
USER ${USER}:${USER}
ENTRYPOINT ["/app/bin/main"]
