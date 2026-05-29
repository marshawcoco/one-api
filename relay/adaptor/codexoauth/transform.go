package codexoauth

import (
	"github.com/songquanpeng/one-api/relay/model"
)

const reasoningEncryptedContent = "reasoning.encrypted_content"

type responsesRequest struct {
	Model             string       `json:"model,omitempty"`
	Input             any          `json:"input"`
	Instructions      string       `json:"instructions"`
	Tools             []model.Tool `json:"tools"`
	ParallelToolCalls bool         `json:"parallel_tool_calls"`
	Store             bool         `json:"store"`
	Include           []string     `json:"include"`
	Stream            bool         `json:"stream"`
	ReasoningEffort   *string      `json:"reasoning_effort,omitempty"`
	Metadata          any          `json:"metadata,omitempty"`
	ServiceTier       *string      `json:"service_tier,omitempty"`
	ToolChoice        any          `json:"tool_choice,omitempty"`
	User              string       `json:"user,omitempty"`
}

func convertToResponsesRequest(request *model.GeneralOpenAIRequest) responsesRequest {
	parallelToolCalls := true
	if request.ParallelTooCalls != nil {
		parallelToolCalls = *request.ParallelTooCalls
	}

	return responsesRequest{
		Model:             request.Model,
		Input:             responsesInput(request),
		Instructions:      request.Instruction,
		Tools:             request.Tools,
		ParallelToolCalls: parallelToolCalls,
		Store:             false,
		Include:           []string{reasoningEncryptedContent},
		Stream:            true,
		ReasoningEffort:   request.ReasoningEffort,
		Metadata:          request.Metadata,
		ServiceTier:       request.ServiceTier,
		ToolChoice:        request.ToolChoice,
		User:              request.User,
	}
}

func responsesInput(request *model.GeneralOpenAIRequest) any {
	if request.Input != nil {
		return request.Input
	}
	if len(request.Messages) == 0 {
		return ""
	}

	input := make([]map[string]any, 0, len(request.Messages))
	for _, message := range request.Messages {
		item := map[string]any{
			"role":    message.Role,
			"content": message.Content,
		}
		if message.ToolCallId != "" {
			item["tool_call_id"] = message.ToolCallId
		}
		if len(message.ToolCalls) > 0 {
			item["tool_calls"] = message.ToolCalls
		}
		input = append(input, item)
	}
	return input
}
