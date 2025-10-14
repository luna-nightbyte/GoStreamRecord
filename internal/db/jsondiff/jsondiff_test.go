package jsondiff

import (
	"os"
	"testing"
)

func TestFiles(t *testing.T) {
	opts := DefaultJSONOptions()
	opts.SkipValueDiff = true

	a, _ := os.ReadFile("a.json")
	b, _ := os.ReadFile("b.json")
	tmp, diff := Compare(a, b, &opts)
	if diff != "" {
		t.Log(diff)
	}
	t.Log(tmp)
}
