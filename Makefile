.DEFAULT_GOAL := dev

.PHONY: dev
dev: ## dev build
dev: clean install generate vet fmt lint test mod-tidy gh-action

.PHONY: gh-action
gh-action: dev dockerize

.PHONY: dockerize
dockerize:
	cp -R laraboot/ actions/commander/entrypoint
	cp -R cmd/ actions/commander/entrypoint
	cp go.mod actions/commander/entrypoint
	cp go.sum actions/commander/entrypoint

.PHONY: ci
ci: ## CI build
ci: dev

.PHONY: clean
clean: ## remove files created during build pipeline
	$(call print-target)
	rm -rf dist
	rm -f coverage.*
	chmod -R +x scripts

.PHONY: tests
tests: ## remove files created during build pipeline
	go get github.com/markbates/pkger/cmd/pkger
	wget -O wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh
	sudo mv ./wait-for-it.sh /usr/bin/wait-for-it
	sudo chmod +x /usr/bin/wait-for-it
	sudo chmod +x -R ./test
	# Prevent permission issues
	sudo chmod 666 /var/run/docker.sock
	./test/buildpack-test.sh
	./test/laraboot-test.sh

.PHONY: install
install: ## go install tools
	$(call print-target)
	cd tools && go install $(shell cd tools && go list -f '{{ join .Imports " " }}' -tags=tools)

.PHONY: generate
generate: ## go generate
	$(call print-target)
	go generate ./...

.PHONY: vet
vet: ## go vet
	$(call print-target)
	go vet ./...

.PHONY: fmt
fmt: ## go fmt
	$(call print-target)
	go fmt ./...

.PHONY: lint
lint: ## golangci-lint
	$(call print-target)
	golangci-lint run

.PHONY: test
test: ## go test with race detector and code covarage
	$(call print-target)
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: mod-tidy
mod-tidy: ## go mod tidy
	$(call print-target)
	go mod tidy
	cd tools && go mod tidy

.PHONY: diff
diff: ## git diff
	$(call print-target)
	git diff --exit-code
	RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi

.PHONY: build
build: ## goreleaser --snapshot --skip-publish --rm-dist
build: install
	$(call print-target)
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: release
release: ## goreleaser --rm-dist
release: install
	$(call print-target)
	goreleaser --rm-dist

.PHONY: run
run: ## go run
	@go run -race .

.PHONY: go-clean
go-clean: ## go clean build, test and modules caches
	$(call print-target)
	go clean -r -i -cache -testcache -modcache

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
