package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/thekhanj/digikala-sdk/schemas/parser"
)

const (
	CODE_SUCCESS int = iota
	CODE_GENERAL_ERROR
)

func parseFile(filePath string) error {
	p, err := parser.NewParser()
	if err != nil {
		return err
	}

	ret, err := p.ParseFile(filePath)
	if err != nil {
		return err
	}
	p.SetSchemas(ret)

	b, err := json.MarshalIndent(ret, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: Not enough arguments")
		os.Exit(CODE_GENERAL_ERROR)
		return
	}

	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "error: Extra argument:", os.Args[2])
		os.Exit(CODE_GENERAL_ERROR)
		return
	}

	filePath := os.Args[1]

	err := parseFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(CODE_GENERAL_ERROR)
	}
}
