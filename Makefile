CONFIG_SCHEMA = config-schema.json

GO_CONFIG = cli/config/dto.go
GO_SRC_FILES = $(shell find cli -name '*.go') \
							 $(wildcard *.go) $(GO_CONFIG)

SCHEMAS = $(shell find schemas -name '*.json')
API_SCHEMAS = $(shell find api -name '*.json')
API_GO_SRC_FILES = $(shell find api -name '*.go')

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
	go build -o $@ ./api

openapi.json: api-gen $(SCHEMAS) $(API_SCHEMAS)
	./api-gen api/openapi.json >$@

$(GO_CONFIG): $(CONFIG_SCHEMA)
	dir=$$(mktemp -d) && \
	cp $< $$dir/config && \
	go-jsonschema -p config $$dir/config >$@; \
	rm $$dir -r
