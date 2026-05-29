package codexoauth

import (
	"testing"

	"github.com/songquanpeng/one-api/relay/meta"
)

func TestGetRequestURLUsesDefaultBaseURL(t *testing.T) {
	url, err := (&Adaptor{}).GetRequestURL(&meta.Meta{})
	if err != nil {
		t.Fatal(err)
	}
	if url != defaultBaseURL+responsesPath {
		t.Fatalf("url mismatch: %s", url)
	}
}

func TestGetRequestURLUsesConfiguredBaseURL(t *testing.T) {
	url, err := (&Adaptor{}).GetRequestURL(&meta.Meta{BaseURL: "https://example.com/codex/"})
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://example.com/codex/responses" {
		t.Fatalf("url mismatch: %s", url)
	}
}
