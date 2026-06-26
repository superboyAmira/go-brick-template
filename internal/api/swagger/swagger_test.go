package swagger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEmbeddedOpenAPIMatchesContract(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	srcPath := filepath.Join(root, "docs", "contracts", "openapi.yaml")
	src, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("read contract: %v (run from module root)", err)
	}
	if len(openAPIYAML) == 0 {
		t.Fatal("embedded openapi is empty — run: make swagger-gen")
	}
	if string(src) != string(openAPIYAML) {
		t.Fatalf("embedded openapi out of sync with %s — run: make swagger-gen", srcPath)
	}
}
