clean:
	@echo ">> cleaning..."
	@rm -rf build/

build/json-pointer: clean
	@echo ">> building..."
	@mkdir -p build
	@cd cli/json-pointer && go build -o ../../build/json-pointer

build/json-patch: clean
	@echo ">> building..."
	@mkdir -p build
	@cd cli/json-patch && go build -o ../../build/json-patch

.PHONY: test
test:
	go test -v ./...

.PHONY: benchmark
benchmark:
	go test -v -run=XXX -bench=.

test_pointer: build/json-pointer
	npx -p json-joy@2.4.0 json-pointer-test ./build/json-pointer

test_patch: build/json-patch
	npx -p json-joy@2.4.0 json-patch-test ./build/json-patch
