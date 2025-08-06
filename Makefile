.PHONY: test-goreadme
test-goreadme:
	@CI=false go run cmd/goreadme/main.go -title='goreadme' > README.md

.PHONY: goreadme
goreadme:
	@CI=false go run github.com/posener/goreadme/cmd/goreadme@v1.4.2 -title='goreadme' > README.md

PKG_LIST := $(shell go list ./...)

.PHONY: build
build:
	go build -o bin/goreadme main.go

.PHONY: run
run: build-local
	./bin/goreadme

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: staticcheck
staticcheck:
	staticcheck ${PKG_LIST}

.PHONY: fmt
fmt:
	go fmt ${PKG_LIST}

.PHONY: vet
vet:
	go vet ${PKG_LIST}

.PHONY: unit
unit:
	go test -v ${PKG_LIST} -count=1 -timeout=10s

.PHONY: race
race:
	go test -race -short ${PKG_LIST}  -count=1 -timeout=10s

.PHONY: gosec
gosec:
	gosec ${PKG_LIST}

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out

# Just for using locally to see the coverage report via html format.
.PHONY: coverage-html
coverage-html:
	mkdir -p "coverage"
	go test -covermode=count -coverprofile "coverage/coverage.cov" ${PKG_LIST}
	go tool cover -func="coverage/coverage.cov"
	go tool cover -html="coverage/coverage.cov" -o "coverage.html"
	open "coverage.html"
	rm -rf "coverage"

.PHONY: gotestsum
gotestsum:
	gotestsum --watch -- --count=1 --timeout=5s
