GO_VERSION := "1.16"
VERSION_DEVELOP := "develop"
GIT_COMMIT := $(shell git log -1 --pretty=format:%h)

develop-images:
	docker build \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg VERSION=${VERSION_DEVELOP} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		-f ninhydrin.Dockerfile -t lostinsoba/ninhydrin:${VERSION_DEVELOP} .

develop-compose: develop-images
	docker-compose -f develop/compose/storage.yml -f develop/compose/ninhydrin.yml up
