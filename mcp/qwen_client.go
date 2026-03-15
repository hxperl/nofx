package mcp

import (
	"net/http"
	"strings"
)

const (
	ProviderQwen       = "qwen"
	DefaultQwenBaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	DefaultQwenModel   = "qwen3-max"
)

type QwenClient struct {
	*Client
}

// NewQwenClient creates Qwen client (backward compatible)
//
// Deprecated: Recommend using NewQwenClientWithOptions for better flexibility
func NewQwenClient() AIClient {
	return NewQwenClientWithOptions()
}

// NewQwenClientWithOptions creates Qwen client (supports options pattern)
//
// Usage examples:
//   // Basic usage
//   client := mcp.NewQwenClientWithOptions()
//
//   // Custom configuration
//   client := mcp.NewQwenClientWithOptions(
//       mcp.WithAPIKey("sk-xxx"),
//       mcp.WithLogger(customLogger),
//       mcp.WithTimeout(60*time.Second),
//   )
func NewQwenClientWithOptions(opts ...ClientOption) AIClient {
	// 1. Create Qwen preset options
	qwenOpts := []ClientOption{
		WithProvider(ProviderQwen),
		WithModel(DefaultQwenModel),
		WithBaseURL(DefaultQwenBaseURL),
	}

	// 2. Merge user options (user options have higher priority)
	allOpts := append(qwenOpts, opts...)

	// 3. Create base client
	baseClient := NewClient(allOpts...).(*Client)

	// 4. Create Qwen client
	qwenClient := &QwenClient{
		Client: baseClient,
	}

	// 5. Set hooks to point to QwenClient (implement dynamic dispatch)
	baseClient.hooks = qwenClient

	return qwenClient
}

func (qwenClient *QwenClient) SetAPIKey(apiKey string, customURL string, customModel string) {
	qwenClient.APIKey = apiKey

	if len(apiKey) > 8 {
		qwenClient.logger.Infof("🔧 [MCP] Qwen API Key: %s...%s", apiKey[:4], apiKey[len(apiKey)-4:])
	}
	if customURL != "" {
		qwenClient.BaseURL = customURL
		qwenClient.logger.Infof("🔧 [MCP] Qwen using custom BaseURL: %s", customURL)
	} else {
		qwenClient.logger.Infof("🔧 [MCP] Qwen using default BaseURL: %s", qwenClient.BaseURL)
	}
	if customModel != "" {
		qwenClient.Model = customModel
		qwenClient.logger.Infof("🔧 [MCP] Qwen using custom Model: %s", customModel)
	} else {
		qwenClient.logger.Infof("🔧 [MCP] Qwen using default Model: %s", qwenClient.Model)
	}
}

func (qwenClient *QwenClient) setAuthHeader(reqHeaders http.Header) {
	qwenClient.Client.setAuthHeader(reqHeaders)
}

// buildMCPRequestBody overrides the base client to apply model-specific defaults
func (qwenClient *QwenClient) buildMCPRequestBody(systemPrompt, userPrompt string) map[string]any {
	requestBody := qwenClient.Client.buildMCPRequestBody(systemPrompt, userPrompt)

	// Model-specific config for qwen3.5-122b
	if strings.Contains(strings.ToLower(qwenClient.Model), "qwen3.5-122b") {
		requestBody["temperature"] = 1.0
		requestBody["top_p"] = 1.0
		requestBody["top_k"] = 40
		requestBody["min_p"] = 0.0
		requestBody["presence_penalty"] = 2.0
		requestBody["repetition_penalty"] = 1.0
		qwenClient.logger.Infof("🔧 [MCP] Applied qwen3.5-122b model-specific config: temperature=1.0, top_p=1.0, top_k=40, min_p=0.0, presence_penalty=2.0, repetition_penalty=1.0")
	}

	return requestBody
}
