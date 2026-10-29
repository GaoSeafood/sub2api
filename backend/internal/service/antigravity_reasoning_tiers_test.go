package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAntigravityReasoningTiersFromClaudeBody(t *testing.T) {
	for _, tc := range []struct {
		name, body, want string
	}{
		{"default", `{}`, "high"},
		{"disabled overrides effort", `{"thinking":{"type":"disabled"},"output_config":{"effort":"high"}}`, "low"},
		{"effort overrides budget", `{"thinking":{"type":"enabled","budget_tokens":24576},"output_config":{"effort":"low"}}`, "low"},
		{"medium effort", `{"output_config":{"effort":"medium"}}`, "medium"},
		{"max effort", `{"output_config":{"effort":"max"}}`, "high"},
		{"medium budget", `{"thinking":{"type":"enabled","budget_tokens":4000}}`, "medium"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, geminiThinkingLevelFromClaudeBody([]byte(tc.body)))
		})
	}
}

func TestAntigravityReasoningTiersRespectExplicitMappings(t *testing.T) {
	for _, model := range []string{"claude-opus-5-5", "claude-sonnet-5-5", "gemini-4-argon"} {
		t.Run(model, func(t *testing.T) {
			mapping := map[string]string{
				model: model + "-high", model + "-low": model + "-low",
				model + "-medium": model + "-medium", model + "-high": model + "-high",
			}
			gateway := &AntigravityGatewayService{}
			account := newAntigravityAccountWithMapping(mapping)
			require.Equal(t, model+"-low", gateway.getMappedModelForThinkingLevel(account, model, "low"))
			require.Equal(t, model+"-medium", gateway.getMappedModelForThinkingLevel(account, model, "medium"))
			require.Equal(t, model+"-high", gateway.getMappedModelForThinkingLevel(account, model, ""))
			require.Equal(t, model+"-high", gateway.getMappedModelForThinkingLevel(account, model+"-high", "low"))
			mapping[model] = model + "-medium"
			account = newAntigravityAccountWithMapping(mapping)
			require.Equal(t, model+"-medium", gateway.getMappedModelForThinkingLevel(account, model, "low"))
			mapping[model] = "custom-alias"
			account = newAntigravityAccountWithMapping(mapping)
			require.Equal(t, "custom-alias", gateway.getMappedModelForThinkingLevel(account, model, "low"))
		})
	}
}
