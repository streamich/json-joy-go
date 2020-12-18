COMMIT_ID ?= $(shell git rev-parse --short HEAD)
VERSION ?= $(shell cat VERSION)

clean:
	@echo ">> cleaning..."
	@rm -rf build/

build/json-pointer: clean build/json-pointer
	@echo ">> building..."
	@echo "Commit: $(COMMIT_ID)"
	@echo "Version: $(VERSION)"
	@ mkdir -p build
	@cd cli && go build -ldflags "-X main.Version=$(VERSION) -X main.CommitId=$(COMMIT_ID)" -o ../build/json-pointer

.PHONY: test
test:
	go test -v ./...

test_pointer: build/json-pointer
	npx -p json-joy@2.3.5 json-pointer-test ./build/json-pointer
