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

test_pointer: build/json-pointer
	npx -p json-joy@2.3.5 json-pointer-test ./build/json-pointer
