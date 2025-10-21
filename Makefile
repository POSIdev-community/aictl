LOCAL_BIN := $(shell pwd)/bin
VERSION := $(shell cat VERSION)
BUILD_OPTIONS=-ldflags="-X 'github.com/POSIdev-community/aictl/pkg/version.version=$(VERSION)' -s -w" -trimpath

export GOBIN=$(LOCAL_BIN)
export PATH:=$(LOCAL_BIN):${PATH}

all: generate test build doc
quick: generate build
pre-commit: generate test doc

.ensure_bin:
	@mkdir -p ${LOCAL_BIN}

.install_mockery:
	@echo -n "⇒ Installing mockery... "
	@go install github.com/vektra/mockery/v3@v3.5.1 >/dev/null 2>&1
	@echo "$$(mockery version) ✅"

install_tools: .install_mockery

generate: install_tools
	@echo -n "⇒ Generating mocks... "
	@mockery --log-level error
	@echo "✅"
	@echo -n "⇒ Running go generate... "
	@go generate ./...
	@echo "✅"

.PHONY: build
build:
	@echo -n "⇒ Building with $(BUILD_OPTIONS)... "
	@go build $(BUILD_OPTIONS) -o bin/aictl cmd/run/main.go
	@echo "✅"

install:
	@echo -n "⇒ Copy aictl to /usr/bin/aictl..."
	@sudo cp bin/aictl /usr/bin/aictl
	@echo "✅"

bash:
	@echo -n "⇒ Add bash completion..."
	@bin/aictl completion bash > /etc/bash_completion.d/aictl
	@echo "✅"

docker:
	@docker build -t "aictl:$(VERSION)" .

docker-file:
	@docker save -o bin/aictl_$(VERSION).tar aictl:$(VERSION)

.PHONY: test
test:
	@echo "⇒ Running tests..."
	@go test -race ./...
	@echo "⇒ Tests ✅"

clean:
	@echo -n "⇒ Cleaning... "
	@rm -rf ./bin
	@echo "✅"

.PHONY: doc
doc:
	@echo -n "⇒ Generate documentation... "
	@go run ./cmd/doc/generate_doc.go
	@git add ./doc/*
	@echo "✅"
