# build docker images
GO_VERSION := "1.19"
GIT_COMMIT := $(shell git log -1 --pretty=format:%h)

go-version:
	@echo ${GO_VERSION}

# build develop image
VERSION_DEVELOP := "develop"

develop-image:
	docker build \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg VERSION=${VERSION_DEVELOP} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		-f ninhydrin.Dockerfile -t lostinsoba/ninhydrin:${VERSION_DEVELOP} .

# launch development environment
DEVELOP_STORAGE := "postgres"

develop-compose: develop-image
	docker-compose \
		-f develop/compose/monitoring.yml \
		-f develop/compose/network.yml \
		-f develop/compose/storage.${DEVELOP_STORAGE}.yml \
		-f develop/compose/ninhydrin.yml up

develop: develop-compose

# publish nightly
VERSION_NIGHTLY := "nightly"

nightly:
	docker build \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg VERSION=${VERSION_NIGHTLY} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		-f ninhydrin.Dockerfile -t ghcr.io/lostinsoba/ninhydrin:${VERSION_NIGHTLY} .
	docker push ghcr.io/lostinsoba/ninhydrin:${VERSION_NIGHTLY}

# run documentation service
DOC_ENV := ".venv"

docs-web:
	python3 -m venv ${DOC_ENV}
	. ${DOC_ENV}/bin/activate && ${DOC_ENV}/bin/pip install -r ./docs/requirements.txt
	. ${DOC_ENV}/bin/activate && mkdocs serve -f ./docs/mkdocs.yml

# run linter
LINTER_VERSION := "1.50.1"

linter-version:
	@echo ${LINTER_VERSION}

lint:
	docker run \
		--rm -v $$(pwd):/app -w /app \
		golangci/golangci-lint:v${LINTER_VERSION} golangci-lint run -v
