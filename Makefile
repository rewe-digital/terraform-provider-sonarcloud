TESTARGS?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=rewe-digital.com
NAMESPACE=platform
NAME=sonarcloud
VERSION?=0.1.0-local
OS_ARCH?=linux_amd64
BINARY_GLOB=./dist/terraform-provider-${NAME}_${OS_ARCH}/terraform*

default: install

update-services-json:
	curl https://sonarcloud.io/api/webservices/list | jq -r . > gen/services.json

gen:
	go generate

build:
	GORELEASER_CURRENT_TAG=$(VERSION) goreleaser build --snapshot --rm-dist

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv $(BINARY_GLOB) ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

fmt:
	go fmt ./sonarcloud ./gen

.PHONY: gen
