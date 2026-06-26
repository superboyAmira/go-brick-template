// Command openapi-validate checks OpenAPI 3.x documents (contract-first source of truth).
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: openapi-validate [--strict] <path-to-openapi.yaml> [more files...]")
		os.Exit(2)
	}
	args := os.Args[1:]
	strict := false
	if len(args) > 0 && args[0] == "--strict" {
		strict = true
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: openapi-validate [--strict] <path-to-openapi.yaml> [more files...]")
		os.Exit(2)
	}
	ctx := context.Background()
	var failed bool
	for _, path := range args {
		if err := validateFile(ctx, path, strict); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			failed = true
		} else {
			fmt.Printf("%s: OK\n", path)
		}
	}
	if failed {
		os.Exit(1)
	}
}

func validateFile(ctx context.Context, path string, strict bool) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = false
	doc, err := loader.LoadFromData(raw)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	if strict {
		if err := doc.Validate(ctx); err != nil {
			return fmt.Errorf("openapi strict validate: %w", err)
		}
	}
	if doc.Info == nil || strings.TrimSpace(doc.Info.Title) == "" {
		return fmt.Errorf("missing info.title")
	}
	if doc.Paths == nil || doc.Paths.Len() == 0 {
		return fmt.Errorf("no paths defined")
	}
	if len(doc.Servers) == 0 {
		return fmt.Errorf("no servers defined")
	}
	return nil
}
