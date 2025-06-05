MAIN_PACKAGE := ./cmd/cli
BINARY_NAME := envsync

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## build: build the production code
.PHONY: build
build:
	@echo "Building binary..."
	@go build -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

## install: install the binary
.PHONY: install
install:
	@sudo cp bin/$(BINARY_NAME) /usr/local/bin
	@echo "Installed binary ✅"

## run: run the production code
.PHONY: run
run:
	@echo "Running binary..."
	@go run $(MAIN_PACKAGE)

## dev: run the code development environment
.PHONY: dev
dev:
	@echo "Running development environment..."
	@go run $(MAIN_PACKAGE)

## watch: run the application with reloading on file changes
.PHONY: watch
watch:
	@air \
		--build.cmd "make build" --build.bin "bin/${BINARY_NAME}" --build.delay "100" \
        --build.exclude_dir "" \
        --build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico, pkl" \
        --misc.clean_on_exit "true"

## update: updates the packages and tidy the modfile
.PHONY: watch
update:
	@go get -u ./...
	@go mod tidy -v


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	@echo "Tidying up..."
	@go fmt ./...
	@go mod tidy -v

## lint: run linter
.PHONY: lint
lint:
	@echo "Linting..."
	@golangci-lint run

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
