# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/test-runner-service ./cmd/server

# Final image: alpine with docker CLI and kubectl installed so the service can drive docker or k8s
FROM alpine:3.18
RUN apk add --no-cache ca-certificates bash curl docker-cli
# install kubectl
ARG KUBECTL_VERSION="v1.27.3"
RUN curl -L -o /usr/local/bin/kubectl https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    && chmod +x /usr/local/bin/kubectl

COPY --from=builder /out/test-runner-service /usr/local/bin/test-runner-service
COPY data /app/data
COPY logo.png /app/logo.png
WORKDIR /app
EXPOSE 9001
ENTRYPOINT ["/usr/local/bin/test-runner-service"]
