.PHONY: all
export GO111MODULE=on

APP=diceDB
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

help:
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done


compile:
	go build -o $(APP_EXECUTABLE) cmd/*.go

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

lint:
	@for p in $(ALL_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

build: # complies format vet and lint your code
build: compile fmt vet lint

set-packages-to-test:
	$(eval PACKAGES_TO_TEST := $(shell go list ./...))

test-cover-html:
	mkdir -p out/
	go test -covermode=count  -coverprofile=cover.out.tmp ${PACKAGES_TO_TEST}
	cat cover.out.tmp | grep -v "_mock.go"  > coverage-all.out
	rm cover.out.tmp
	@go tool cover -html=coverage-all.out -o out/coverage.html
	@go tool cover -func=coverage-all.out

test: # run test's for your code and generate out a test report
test: set-packages-to-test test-cover-html

test-cover:
	go get -u github.com/jokeyrhyme/go-coverage-threshold/cmd/go-coverage-threshold
	ENVIRONMENT=test go-coverage-threshold

run: # run's your go code from source
run:
	go run cmd/main.go -host=127.0.0.1 -port=7379

clean:
	rm -rf out/