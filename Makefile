BUI_LIB?=github.com/ggrrrr/bui_lib_http
PROJECT?=github.com/ggrrrr/bui_api_login
APP?=api
PORT?=8000

RELEASE?=0.0.1
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_OS?=$(shell uname)
CONTAINER_IMAGE?=${APP}

GOOS?=linux
GOARCH?=amd64

clean:
	rm -f ${APP}

build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w \
		-X ${BUI_LIB}/build.Release=${RELEASE} \
		-X ${BUI_LIB}/build.Commit=${COMMIT} \
		-X ${BUI_LIB}/build.BuildTime=${BUILD_TIME}" \
		-o ${APP}

buildmac: clean
	CGO_ENABLED=0 go build \
		-ldflags "-s -w \
		-X ${BUI_LIB}/build.Release=${RELEASE} \
		-X ${BUI_LIB}/build.Commit=${COMMIT} \
		-X ${BUI_LIB}/build.BuildTime=${BUILD_TIME} \
		-X ${BUI_LIB}/build.BuildOs=${BUILD_OS}" \
		-o ${APP}

container: build
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

run: container
	docker stop $(APP):$(RELEASE) || true && docker rm $(APP):$(RELEASE) || true
	docker run --name ${APP} -p ${PORT}:${PORT} --rm \
		-e "PORT=${PORT}" \
		$(APP):$(RELEASE)

test:
	go test -v -race ./...

push: container
	docker push $(CONTAINER_IMAGE):$(RELEASE)

minikube: push
	for t in $(shell find ./kubernetes/advent -type f -name "*.yaml"); do \
        cat $$t | \
        	gsed -E "s/\{\{(\s*)\.Release(\s*)\}\}/$(RELEASE)/g" | \
        	gsed -E "s/\{\{(\s*)\.ServiceName(\s*)\}\}/$(APP)/g"; \
        echo ---; \
    done > tmp.yaml
	kubectl apply -f tmp.yaml