FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye

# [Optional] Uncomment the next lines to use go get to install anything else you need
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y upgrade && apt-get -y install protobuf-compiler libprotobuf-dev
