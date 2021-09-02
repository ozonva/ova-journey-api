GOBIN?=$(GOPATH)/bin
LOCAL_BIN:=$(CURDIR)/bin
THIRD_PARTY:=$(CURDIR)/third_party
PROTO_API:=$(CURDIR)/api

GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.17","$(shell printf "$(GO_VERSION_SHORT)\n1.17" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.17. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on
export GOPROXY=https://proxy.golang.org|direct

PGV_VERSION:="v0.6.1"
GOOGLEAPIS_VERSION="master"
BUF_VERSION:="v0.51.0"

.PHONY: vendor-proto
vendor-proto:
	@[ -f $(THIRD_PARTY)/validate/validate.proto ] || (mkdir -p $(THIRD_PARTY)/validate/ && curl -sSL0 https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/$(PGV_VERSION)/validate/validate.proto -o $(THIRD_PARTY)/validate/validate.proto)
	@[ -f $(THIRD_PARTY)/google/api/http.proto ] || (mkdir -p $(THIRD_PARTY)/google/api/ && curl -sSL0 https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_VERSION)/google/api/http.proto -o $(THIRD_PARTY)/google/api/http.proto)
	@[ -f $(THIRD_PARTY)/google/api/annotations.proto ] || (mkdir -p $(THIRD_PARTY)/google/api/ && curl -sSL0 https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_VERSION)/google/api/annotations.proto -o $(THIRD_PARTY)/google/api/annotations.proto)

buf.work.yaml:
	@echo "version: v1\ndirectories:\n  - api\n  - third_party\n" > $(CURDIR)/buf.work.yaml
buf.gen.yaml:
	@echo "version: v1\nplugins:\n  - name: go\n    out: .\n    opt: module=github.com/ozonva/ova-journey-api\n  - name: go-grpc\n    out: .\n    opt: module=github.com/ozonva/ova-journey-api\n  - name: validate\n    out: .\n    opt: lang=go,module=github.com/ozonva/ova-journey-api\n  - name: grpc-gateway\n    out: .\n    opt: logtostderr=true,module=github.com/ozonva/ova-journey-api\n  - name: openapiv2\n    out: swagger\n    opt: allow_merge=true,merge_file_name=api" > $(CURDIR)/buf.gen.yaml

.PHONY: generate
generate: vendor-proto buf.work.yaml buf.gen.yaml
	@command -v buf 2>&1 > /dev/null || (mkdir -p $(GOBIN) && curl -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(shell uname -s)-$(shell uname -m) -o $(GOBIN)/buf && chmod +x $(GOBIN)/buf)
	@[ -f $(PROTO_API)/buf.yaml ] || buf mod init --doc -o $(PROTO_API)
	buf generate $(PROTO_API)
	go generate ./...

.PHONY: deps
deps:
	@[ -f go.mod ] || go mod init github.com/ozonva/ova-journey-api
	find . -name go.mod -exec bash -c 'pushd "$${1%go.mod}" && go mod tidy && popd' _ {} \;

.PHONY: bin-deps
bin-deps:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go install github.com/envoyproxy/protoc-gen-validate@$(PGV_VERSION)

.PHONY: run
run:
	go run cmd/ova-journey-api/main.go

.PHONY: test
test:
	go test -race ./...
	go test ./...

.PHONY: build
build: deps
	CGO_ENABLED=0 go build -o $(LOCAL_BIN)/ova-journey-api cmd/ova-journey-api/main.go

.PHONY: clean
clean:
	rm -f $(CURDIR)/buf.work.yaml
	rm -f $(CURDIR)/buf.gen.yaml
	rm -f swagger/api.swagger.json
	rm -rf $(THIRD_PARTY)
	rm -rf $(LOCAL_BIN)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: docker-build
docker-build:
	docker-compose build

.PHONY: all
all: clean bin-deps deps generate lint test build