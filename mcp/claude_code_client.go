package mcp

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	ProviderClaudeCode       = "claude-code"
	DefaultClaudeCodeModel   = "claude-sonnet-4-6"
	DefaultClaudeCodeTimeout = 300 * time.Second
)

// ClaudeCodeClient uses the Claude Code CLI (claude) to make AI calls.
// No API key required — it uses the locally authenticated CLI session.
type ClaudeCodeClient struct {
	model   string
	timeout time.Duration
	logger  Logger
}

// NewClaudeCodeClient creates a new Claude Code CLI client.
func NewClaudeCodeClient() AIClient {
	return &ClaudeCodeClient{
		model:   DefaultClaudeCodeModel,
		timeout: DefaultClaudeCodeTimeout,
		logger:  &noopLogger{},
	}
}

func (c *ClaudeCodeClient) SetAPIKey(apiKey string, customURL string, customModel string) {
	if customModel != "" {
		c.model = customModel
		c.logger.Infof("🔧 [MCP] Claude Code CLI using custom model: %s", customModel)
	} else {
		c.logger.Infof("🔧 [MCP] Claude Code CLI using default model: %s", c.model)
	}
}

func (c *ClaudeCodeClient) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *ClaudeCodeClient) CallWithMessages(systemPrompt, userPrompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	args := []string{
		"-p", userPrompt,
		"--output-format", "text",
		"--allowedTools", "",
	}

	if systemPrompt != "" {
		args = append(args, "--system-prompt", systemPrompt)
	}

	if c.model != "" {
		args = append(args, "--model", c.model)
	}

	c.logger.Infof("🤖 [MCP] Claude Code CLI executing with timeout %v", c.timeout)

	cmd := exec.CommandContext(ctx, "claude", args...)

	// Clear CLAUDECODE env var to prevent "nested session" detection
	cmd.Env = filterEnv(os.Environ(), "CLAUDECODE")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	if err := cmd.Run(); err != nil {
		elapsed := time.Since(startTime)
		stderrStr := strings.TrimSpace(stderr.String())
		stdoutStr := strings.TrimSpace(stdout.String())
		c.logger.Errorf("🤖 [MCP] Claude Code CLI failed after %v, stdout(%d bytes): %.200s, stderr: %s", elapsed, len(stdoutStr), stdoutStr, stderrStr)
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("claude CLI timed out after %v: %s", c.timeout, stderrStr)
		}
		return "", fmt.Errorf("claude CLI error: %w, stderr: %s", err, stderrStr)
	}

	elapsed := time.Since(startTime)
	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", fmt.Errorf("claude CLI returned empty response after %v", elapsed)
	}

	c.logger.Infof("🤖 [MCP] Claude Code CLI responded in %v (%d bytes)", elapsed, len(result))
	return result, nil
}

func (c *ClaudeCodeClient) CallWithRequest(req *Request) (string, error) {
	var systemPrompt, userPrompt string

	for _, msg := range req.Messages {
		switch msg.Role {
		case "system":
			systemPrompt = msg.Content
		case "user":
			if userPrompt != "" {
				userPrompt += "\n\n"
			}
			userPrompt += msg.Content
		case "assistant":
			if userPrompt != "" {
				userPrompt += "\n\n"
			}
			userPrompt += fmt.Sprintf("[Previous assistant response]: %s", msg.Content)
		}
	}

	if userPrompt == "" {
		return "", fmt.Errorf("no user message found in request")
	}

	return c.CallWithMessages(systemPrompt, userPrompt)
}

// filterEnv returns a copy of env with the specified key removed.
func filterEnv(env []string, key string) []string {
	prefix := key + "="
	result := make([]string, 0, len(env))
	for _, e := range env {
		if !strings.HasPrefix(e, prefix) {
			result = append(result, e)
		}
	}
	return result
}
