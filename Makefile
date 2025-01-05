.PHONY: help clean coverage pkgsite report test vuln chores

help: ## list available targets
	@# Shamelessly stolen from Gomega's Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

clean: ## cleans up build and testing artefacts
	rm -f coverage.*

test: ## run unit tests
	go test -v -p=1 -race ./...

report: ## run goreportcard-cli on this module
# from ghcr.io/thediveo/devcontainer-features/goreportcard
	goreportcard-cli -v ./..

coverage: ## gathers coverage and updates README badge
# from ghcr.io/thediveo/devcontainer-features/gocover
	gocover
