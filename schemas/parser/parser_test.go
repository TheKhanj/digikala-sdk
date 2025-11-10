package parser

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {
	p, err := NewParser()
	if err != nil {
		t.Fatal(err)
		return
	}

	ret, err := p.ParseFile(
		filepath.Join("../../openapi/openapi.json"),
	)
	if err != nil {
		t.Fatal(err)
		return
	}
	p.SetSchemas(ret)

	b, err := json.Marshal(ret)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(b))
}
