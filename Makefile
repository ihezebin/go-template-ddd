#BRANCH ?= `git rev-parse --abbrev-ref HEAD`
#COMMIT ?= `git rev-parse --short HEAD`
TAG ?= `git describe --tags --always`
PROJECT_ROOT = `pwd`
PROJECT_NAME = go-template-ddd
#阿里云：registry.cn-chengdu.aliyuncs.com
#腾讯云：ccr.ccs.tencentyun.com
DOCKER_REGISTRY ?= $(HEZEBIN_DOCKER_REGISTRY)
DOCKER_USER ?= $(HEZEBIN_DOCKER_USER)
DOCKER_PWD ?= $(HEZEBIN_DOCKER_PWD)
DOCKER_TAG ?= $(DOCKER_REGISTRY)/hezebin/$(PROJECT_NAME):$(TAG)

.PHONY: package
package: tag clean

.PHONY: build
tag: build
	docker login --username=$(DOCKER_USER) --password=$(DOCKER_PWD) $(DOCKER_REGISTRY)
	docker build --platform linux/amd64 --build-arg PROJECT_NAME=$(PROJECT_NAME) --build-arg TAG=$(TAG) -t $(DOCKER_TAG) -f Dockerfile .
	docker push $(DOCKER_TAG)
	echo $(DOCKER_TAG)

.PHONY: init
init:


.PHONY: test
test:
	go test -tags=unit -timeout 30s -short -v `go list ./... | grep -v 'component'`

.PHONY: build
build: init test
	go generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/$(PROJECT_NAME) $(PKG_ROOT)

.PHONY: clean
clean:
	rm -rf ./build;

