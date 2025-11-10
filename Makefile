BUILD_DIR = build

SCHEMAS = $(shell find schemas -name '*.json')
API_SCHEMAS = $(shell find openapi/api -name '*.json')
GO_SRC_FILES = $(shell find . -name '*.go' | grep -v 'api/api.go' | grep -v 'api/schemas.go')

all: openapi api

.PHONY: clean
clean: clean_cli clean_api clean_openapi clean_test
	rm $(BUILD_DIR) -r

$(BUILD_DIR):
	mkdir $@

cli: $(BUILD_DIR)/cli
.PHONY: clean_cli
clean_cli:
	rm $(BUILD_DIR)/cli

$(BUILD_DIR)/cli: $(GO_SRC_FILES)
	go build -o $@ ./cli

api: api/api.go api/schemas.go
.PHONY: clean_api
clean_api:
	rm api/api.go api/schemas.go api/schemas

api/api.go: api $(BUILD_DIR)/openapi.json
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
		-generate "client" \
		-o /dev/stdout \
		-package api $< \
		> $@

api/schemas: $(BUILD_DIR)/openapi.json
	cat $< \
		| jq '{ "type": "object", "definitions": .components.schemas }' \
		| sed 's|/components/schemas|/definitions|' > $@

api/schemas.go: api/schemas
	go run github.com/atombender/go-jsonschema@latest -p api $< > $@

openapi: $(BUILD_DIR)/openapi.json
.PHONY: clean_openapi
clean_openapi:
	rm $(BUILD_DIR)/openapi.json

$(BUILD_DIR)/openapi.json: $(BUILD_DIR)/cli $(SCHEMAS) $(API_SCHEMAS)
	$< openapi/openapi.json >$@

test: test/config/config.go
.PHONY: clean_test
clean_test:
	rm test/config/config.go test/config/config

test/config/config: test/schema.json
	cat $< >$@

test/config/config.go: test/config/config
	go run github.com/atombender/go-jsonschema@latest -p config $< > $@
