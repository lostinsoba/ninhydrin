ARG VERSION=""
ARG GO_VERSION=""
ARG GIT_COMMIT=""

FROM golang:${GO_VERSION}-alpine as builder

RUN apk add --no-cache git

WORKDIR /go/src/lostinsoba/ninhydrin

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /go/src/lostinsoba/ninhydrin

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-s -w -X main.name=api -X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT} " \
    -o build/api lostinsoba/ninhydrin/cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-s -w -X main.name=scheduler -X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT} " \
    -o build/scheduler lostinsoba/ninhydrin/cmd/scheduler

FROM alpine

LABEL "org.opencontainers.image.title"="Ninhydrin"
LABEL "org.opencontainers.image.description"="Distributed task registry"
LABEL "org.opencontainers.image.url"="https://github.com/lostinsoba/ninhydrin"
LABEL "org.opencontainers.image.licenses"="AGPL-3.0"
LABEL "org.opencontainers.image.version"=${VERSION}
LABEL "org.opencontainers.image.revision"=${GIT_COMMIT}

RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN addgroup -S ninhydrin && adduser -S ninhydrin -G ninhydrin
USER ninhydrin

COPY --from=builder /go/src/lostinsoba/ninhydrin/build ./ninhydrin
COPY ninhydrin.yml /etc/ninhydrin/ninhydrin.yml

EXPOSE 8080 8080
