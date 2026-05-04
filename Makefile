# fluent-sdfx project Makefile.
#
# Conventions:
#   make build       — sync docs from tutorial code, vet, test (the main check).
#   make install-hooks — install git hooks that auto-run `make build` on commit.
#   make serve       — local SPA-fallback dev server for the docs site.
#   make screenshots — regenerate every docs figure via f3d (expensive).

.PHONY: build docs-sync vet test fmt-check screenshots serve install-hooks help

build: docs-sync vet test
	@echo "✓ build clean"

# Re-inline tutorial Go source into matching docs/content/*.md code blocks.
docs-sync:
	@./tools/update-tutorial-code.sh

vet:
	@go vet ./...

test:
	@go test ./...

fmt-check:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt issues:"; echo "$$unformatted"; exit 1; \
	fi

# Regenerate every docs/images/<page>-<step>.png via f3d. Expensive — runs
# every tutorial step's `go run` and a headless f3d render. Not part of
# `make build`; invoke manually when tutorial code or output changes
# meaningfully.
screenshots:
	@./tutorial/render-screenshots.sh

# Local SPA-fallback dev server on :7174.
serve:
	@./tools/serve.py

install-hooks:
	@./tools/install-hooks.sh

help:
	@echo "Targets:"
	@echo "  build         — docs-sync + vet + test"
	@echo "  docs-sync     — re-inline tutorial Go source into docs/content/*.md"
	@echo "  vet           — go vet ./..."
	@echo "  test          — go test ./..."
	@echo "  fmt-check     — gofmt -l . (zero output = clean)"
	@echo "  screenshots   — regenerate every docs figure via f3d"
	@echo "  serve         — local SPA-fallback dev server on :7174"
	@echo "  install-hooks — install git pre-commit hook"
