GO ?= go
GOLANGCI_LINT ?= golangci-lint
PNPM ?= pnpm
export CODEX_HOME := $(PWD)/.codex
export CODEX_HOME_PUB := $(PWD)/.codex_pub

.PHONY: lint build test editor-install editor-dev editor-test editor-test-e2e editor-package editor-package-dir

# Lint/check basics: catches unused imports via build and common vet checks
lint:
	@if command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		$(GOLANGCI_LINT) run; \
	else \
		echo "$(GOLANGCI_LINT) not installed; skipping golangci-lint"; \
	fi
	$(GO) vet ./...
	$(GO) build ./...

# Convenience: full build
build:
	$(GO) build ./...

test:
	@if [ -f go.mod ]; then \
		$(GO) test ./...; \
	else \
		echo "go.mod not found; skipping go test"; \
	fi
	@if [ -f packages/bilink-npx/test/platform.test.mjs ]; then \
		node --test packages/bilink-npx/test/platform.test.mjs; \
	fi

editor-install:
	cd apps/editor && $(PNPM) install

editor-dev:
	cd apps/editor && $(PNPM) run dev

editor-test: editor-test-e2e

editor-test-e2e:
	cd apps/editor && $(PNPM) run test:e2e

editor-package:
	cd apps/editor && $(PNPM) run package

editor-package-dir:
	cd apps/editor && $(PNPM) run package:dir

codex-locale:
	@echo "CODEX_HOME=$(CODEX_HOME)"
	@echo "Running codex with CODEX_HOME=$(CODEX_HOME)"
	codex --dangerously-bypass-approvals-and-sandbox

codex-locale-resume:
	@echo "CODEX_HOME=$(CODEX_HOME)"
	@echo "Running codex with CODEX_HOME=$(CODEX_HOME)"
	codex --dangerously-bypass-approvals-and-sandbox resume

codex-locale-pub:
	@echo "CODEX_HOME=$(CODEX_HOME_PUB)"
	@echo "Running codex with CODEX_HOME=$(CODEX_HOME_PUB)"
	codex --dangerously-bypass-approvals-and-sandbox
