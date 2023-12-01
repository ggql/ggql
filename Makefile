# Config

VERSION=$(version)


# Build

.PHONY: FORCE

build: go-build
.PHONY: build

clean: go-clean
.PHONY: clean

lint: go-lint
.PHONY: lint

test: go-test
.PHONY: test


# Non-PHONY targets (real files)

go-build: FORCE
	./scripts/build.sh $(VERSION)

go-clean: FORCE
	./scripts/clean.sh

go-lint: FORCE
	./scripts/lint.sh

go-test: FORCE
	./scripts/test.sh report
