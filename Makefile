CONFIG_SCHEMA = config-schema.json

GO_CONFIG = cli/config/dto.go
GO_SRC_FILES = $(shell find cli -name '*.go') \
							 $(wildcard *.go) $(GO_CONFIG)

SCHEMAS = $(shell find schemas -name '*.json')
API_SCHEMAS = $(shell find api -name '*.json')
API_GO_SRC_FILES = $(shell find api -name '*.go' | grep -v 'api.go' | grep -v 'schemas.go')

all: cli openapi.json

cli: bin/cli

bin/cli: $(GO_SRC_FILES)
	go build -o ./bin/cli ./cli

go-config: $(GO_CONFIG)

.PHONY: clean
clean: clean-go-config

.PHONY: clean-go-config
clean-go-config:
	rm $(GO_CONFIG)

api-gen: $(API_GO_SRC_FILES)
	go build -o $@ ./api/cli

api/api.go: openapi.json
	oapi-codegen -generate "client" \
		-o api/api.go \
		-package api openapi.json

openapi.json: api-gen $(SCHEMAS) $(API_SCHEMAS)
	./api-gen api/openapi.json >$@

api/schemas: openapi.json
	cat $< \
		| jq '{ "type": "object", "definitions": .components.schemas }' \
		| sed 's|/components/schemas|/definitions|' > $@

api/schemas.go: api/schemas
	go run github.com/atombender/go-jsonschema@latest -p api $< > $@

$(GO_CONFIG): $(CONFIG_SCHEMA)
	dir=$$(mktemp -d) && \
	cp $< $$dir/config && \
	go-jsonschema -p config $$dir/config >$@; \
	rm $$dir -r

test/config/config: test/schema.json
	cat $< >$@

test/config/config.go: test/config/config
	go run github.com/atombender/go-jsonschema@latest -p config $< > $@
