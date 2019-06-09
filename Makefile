TARGET ?= logtap
RELEASE ?= 0.1.0
REGISTRY ?= lichuan0620
TAG ?= ${RELEASE}

.PHONY: default check dep lint test build build-linux build-local image push clean

default: dep check image

dep:
	dep ensure -v -update

check: lint test

lint:
	golint cmd/... pkg/...

test:
	go test ./...

build: build-local

build-linux:
	GOOS=linux GOARCH=amd64 go build \
		-ldflags="-X 'github.com/lichuan0620/logtap/cmd/${TARGET}/version.Version=${RELEASE}'" \
		-o bin/${TARGET} ./cmd/${TARGET}/main.go

build-local:
	go build \
		-ldflags="-X 'github.com/lichuan0620/logtap/cmd/${TARGET}/version.Version=${RELEASE}'" \
		-o bin/${TARGET} ./cmd/${TARGET}/main.go

image: build-linux
	docker build -t ${REGISTRY}/${TARGET}:${TAG} -f build/${TARGET}/Dockerfile .

push: image
	docker push ${REGISTRY}/${TARGET}:${TAG}

clean:
	rm -rf bin
