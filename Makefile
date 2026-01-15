LOCAL_BIN := $(shell pwd)/bin
VERSION := $(shell cat VERSION)
BUILD_OPTIONS=-ldflags="-X 'github.com/POSIdev-community/aictl/pkg/version.version=$(VERSION)' -s -w" -trimpath

PLATFORMS := linux_amd64 linux_arm64 darwin_amd64 darwin_arm64 windows_amd64

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

.install_enumer:
	@echo -n "⇒ Installing enumer... "
	@go install github.com/dmarkham/enumer@latest >/dev/null 2>&1
	@echo "✅"

.install_golangci-lint:
	@echo -n "⇒ Installing golangci-lint... "
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0 >/dev/null 2>&1
	@echo "✅"

lint: | .install_golangci-lint
	@echo -n "⇒ Linting... "
	@golangci-lint run
	@echo "✅"

install_tools: .install_mockery .install_enumer

generate: install_tools
	@echo -n "⇒ Generating mocks... "
	@mockery --log-level error
	@echo "✅"
	@echo -n "⇒ Running go generate... "
	@go generate ./...
	@echo "✅"

.PHONY: build
build:
	@echo -n "⇒ Building for local platform with $(BUILD_OPTIONS)... "
	@go build $(BUILD_OPTIONS) -o bin/aictl cmd/run/main.go
	@echo "✅"

# Сборка под все указанные платформы
.PHONY: build-all
build-all: | .ensure_bin
	@echo "⇒ Building for all platforms: $(PLATFORMS)..."
	@for platform in $(PLATFORMS); do \
		goos=$$(echo $$platform | cut -d'_' -f1); \
		goarch=$$(echo $$platform | cut -d'_' -f2); \
		ext=""; \
		if [ "$$goos" = "windows" ]; then ext=".exe"; fi; \
		echo "  → $$goos/$$goarch"; \
		GOOS=$$goos GOARCH=$$goarch go build $(BUILD_OPTIONS) -o "bin/aictl_$$goos-$$goarch$$ext" cmd/run/main.go; \
	done
	@echo "✅ All builds completed."

install:
	@echo -n "⇒ Copy aictl to /usr/bin/aictl..."
	@sudo cp bin/aictl /usr/bin/aictl
	@echo "✅"

bash:
	@echo -n "⇒ Add bash completion..."
	@bin/aictl completion bash | sudo tee /etc/bash_completion.d/aictl >/dev/null
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

.PHONY: check-doc
check-doc:
	@go run ./cmd/doc/generate_doc.go
	@git diff --exit-code -- ./doc/ || (echo "❌ Docs outdated. Run 'make doc' and commit."; exit 1)
