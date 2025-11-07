package parser

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Schema any

func NewParser() (*Parser, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Parser{
		schemas:     make(map[string]Schema),
		schemasPath: "#/components/schemas",
		cwd:         wd,
	}, nil
}

type Parser struct {
	schemas     map[string]Schema
	schemasPath string
	cwd         string
}

func (this *Parser) ParseFile(path string) (Schema, error) {
	oldwd := this.cwd
	file, err := os.Open(filepath.Join(this.cwd, path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	this.cwd = filepath.Dir(filepath.Join(this.cwd, path))
	defer func() { this.cwd = oldwd }()

	return this.Parse(b)
}

func (this *Parser) Parse(b []byte) (Schema, error) {
	var obj any
	err := json.Unmarshal(b, &obj)
	if err != nil {
		return nil, err
	}

	this.normalizeSchema(obj)

	return obj, nil
}

func (this *Parser) SetSchemas(obj Schema) {
	path := this.schemasPath[2:]

	mp, lastField, err := this.getMapFromPath(obj, path, true)
	if err != nil {
		panic(err)
	}

	mp[lastField] = this.schemas
}

func (this *Parser) getMapFromPath(obj any, path string, create bool) (map[string]any, string, error) {
	fields := strings.Split(path, "/")

	for _, field := range fields[:len(fields)-1] {
		mp, ok := obj.(map[string]any)
		if !ok {
			return nil, "", errors.New("expected to be a map")
		}
		if mp[field] == nil {
			if !create {
				return nil, "", errors.New("expected to be a map")
			}
			mp[field] = make(map[string]any)
		}
		obj = mp[field]
	}
	lastField := fields[len(fields)-1]

	mp, ok := obj.(map[string]any)
	if !ok {
		return nil, "", errors.New("expected to be a map")
	}

	return mp, lastField, nil
}

func (this *Parser) handleLocalSchema(schema Schema, path string) {
	obj, lastField, err := this.getMapFromPath(schema, path, false)
	if err != nil {
		return
	}

	schemas := obj[lastField]
	if schemas == nil {
		return
	}

	if mp, ok := schemas.(map[string]any); ok {
		for schemaName, schema := range mp {
			this.normalizeSchema(schema)
			this.schemas[schemaName] = schema
		}

		delete(obj, lastField)
	}
}

func (this *Parser) normalizeSchema(schema Schema) {
	if mp, ok := schema.(map[string]any); ok {
		delete(mp, "$schema")
	}

	this.traverse(
		schema,
		func(mp map[string]any, key string) {
			if key != "$ref" {
				return
			}

			ref, ok := mp[key].(string)
			if !ok {
				panic("Invalid ref type: " + ref)
			}
			if this.isFilePathRef(ref) {
				schemaName := this.getSchemaNameFromFilePathRef(ref)
				mp[key] = this.schemasPath + "/" + schemaName
				parsed, err := this.ParseFile(ref)
				if err != nil {
					panic(err)
				}
				this.schemas[schemaName] = parsed
			} else if this.isLocalRef(ref) {
				schemaName := this.getSchemaNameFromLocalRef(ref)
				mp[key] = this.schemasPath + "/" + schemaName
			}
		},
	)

	this.handleLocalSchema(schema, "definitions")
	this.handleLocalSchema(schema, "components/schemas")
}

func (this *Parser) traverse(obj any, fn func(mp map[string]any, key string)) {
	if arr, ok := obj.([]any); ok {
		for _, el := range arr {
			this.traverse(el, fn)
		}
	} else if mp, ok := obj.(map[string]any); ok {
		for key, val := range mp {
			fn(mp, key)
			this.traverse(val, fn)
		}
	}
}

func (this *Parser) getSchemaNameFromLocalRef(ref string) string {
	return filepath.Base(ref)
}

func (this *Parser) getSchemaNameFromFilePathRef(ref string) string {
	base := filepath.Base(ref)
	ret := ""
	for i := 0; i < len(base); {
		c := base[i]
		if c == '-' || i == 0 {
			for ; i < len(base) && base[i] == '-'; i++ {
			}
			ret += strings.ToUpper(string(base[i]))
		} else if c == '.' {
			break
		} else {
			ret += string(c)
		}
		i++
	}

	return string(ret)
}

func (this *Parser) isLocalRef(ref string) bool {
	return ref[0] == '#'
}

func (this *Parser) isFilePathRef(ref string) bool {
	return ref[0] == '.'
}
