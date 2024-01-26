

.PHONY: clean all init generate generate_mocks

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: generate
	go mod tidy
	go mod vendor

test:
	go test -short -coverprofile coverage.out -v ./...

generate: generated generate_mocks generate_mocks_all

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

INTERFACES_GO_FILES := $(shell find repository -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

generate_mocks_all:
	@echo "Generating mocks for all interfaces"
	@for file in $$(find pkg -name "*.go" | grep -v "_test.go" | grep -v "_mock.go"); do \
		echo "Generating mocks for $$file"; \
		src_file=$$file; \
		src_file_mock=$$(dirname $$file)/$$(basename $$file .go)_mock.go; \
		package_name=$$(basename $$(dirname $$file)); \
		package_name_no_ext=$$(basename $${package_name} .go); \
		mockgen -source=$$src_file -destination=$$src_file_mock -package=$$package_name_no_ext; \
	done