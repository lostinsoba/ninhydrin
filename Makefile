# build docker images
GO_VERSION := "1.16"
VERSION_DEVELOP := "develop"
GIT_COMMIT := $(shell git log -1 --pretty=format:%h)

go-version:
	@echo ${GO_VERSION}

develop-image:
	docker build \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg VERSION=${VERSION_DEVELOP} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		-f ninhydrin.Dockerfile -t lostinsoba/ninhydrin:${VERSION_DEVELOP} .

develop-compose: develop-image
	docker-compose \
		-f develop/compose/monitoring.yml \
		-f develop/compose/network.yml \
		-f develop/compose/storage.yml \
		-f develop/compose/ninhydrin.yml up

develop: develop-compose

# publish to ghcr
develop-ghcr:
	docker build \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg VERSION=${VERSION_DEVELOP} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		-f ninhydrin.Dockerfile -t ghcr.io/lostinsoba/ninhydrin:${VERSION_DEVELOP} .
	docker push ghcr.io/lostinsoba/ninhydrin:${VERSION_DEVELOP}

# run linter
LINTER_VERSION := "1.50.1"

linter-version:
	@echo ${LINTER_VERSION}

lint:
	docker run \
		--rm -v $$(pwd):/app -w /app \
		golangci/golangci-lint:v${LINTER_VERSION} golangci-lint run -v
