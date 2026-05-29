package codexoauth

import (
	"testing"

	"github.com/songquanpeng/one-api/relay/model"
)

func TestConvertToResponsesRequestUsesInputWhenPresent(t *testing.T) {
	input := []any{map[string]any{"role": "user", "content": "hello"}}
	request := &model.GeneralOpenAIRequest{
		Model: "gpt-5-codex",
		Input: input,
	}

	converted := convertToResponsesRequest(request)

	if converted.Model != "gpt-5-codex" {
		t.Fatalf("model mismatch: %s", converted.Model)
	}
	if converted.Input == nil {
		t.Fatal("input should be preserved")
	}
	if converted.Store {
		t.Fatal("store must be false")
	}
	if !converted.Stream {
		t.Fatal("stream must be true")
	}
	if len(converted.Include) != 1 || converted.Include[0] != reasoningEncryptedContent {
		t.Fatalf("include mismatch: %#v", converted.Include)
	}
	if !converted.ParallelToolCalls {
		t.Fatal("parallel_tool_calls should default to true")
	}
}

func TestConvertToResponsesRequestConvertsMessages(t *testing.T) {
	request := &model.GeneralOpenAIRequest{
		Model: "gpt-5-codex",
		Messages: []model.Message{
			{Role: "user", Content: "hello"},
		},
	}

	converted := convertToResponsesRequest(request)
	input, ok := converted.Input.([]map[string]any)
	if !ok {
		t.Fatalf("input type mismatch: %T", converted.Input)
	}
	if len(input) != 1 {
		t.Fatalf("input length mismatch: %d", len(input))
	}
	if input[0]["role"] != "user" || input[0]["content"] != "hello" {
		t.Fatalf("input item mismatch: %#v", input[0])
	}
}

func TestConvertToResponsesRequestRespectsParallelToolCalls(t *testing.T) {
	parallel := false
	request := &model.GeneralOpenAIRequest{ParallelTooCalls: &parallel}

	converted := convertToResponsesRequest(request)

	if converted.ParallelToolCalls {
		t.Fatal("parallel_tool_calls should respect request override")
	}
}
