package agent

import (
	"strings"
	"testing"

	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/llm"
)

func TestEconomyCompactionSummarizesBeforeLargeWindowFills(t *testing.T) {
	provider := &autoTitleProvider{reply: "Key facts and source IDs remain available."}
	opts := agentcore.Options{
		Provider: provider, Compaction: compactionConfig(200_000),
		MaxTokens: 32_768, NonStreaming: true,
	}
	installEconomyCompaction(&opts)
	if opts.Compactor == nil || opts.MaxTokens != 32_768 {
		t.Fatal("normal reply configuration changed")
	}
	msgs := make([]llm.Message, 12)
	for i := range msgs {
		msgs[i] = llm.UserText(strings.Repeat("source evidence ", 700))
	}
	out := opts.Compactor.Pre(t.Context(), msgs, 30_000)
	if len(out) <= len(msgs) {
		t.Fatal("expected early summarization while the 200K model window still has room")
	}
	if provider.req.MaxTokens != 4_096 || provider.req.Thinking != "disabled" {
		t.Fatalf("summary request budget = %d, thinking = %q", provider.req.MaxTokens, provider.req.Thinking)
	}
}
