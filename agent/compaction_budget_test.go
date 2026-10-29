package agent

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestCompressionInputBudgetTracksModelWindow(t *testing.T) {
	for _, tc := range []struct{ window, want int }{
		{32_000, 16_000},
		{128_000, 64_000},
		{200_000, 64_000},
	} {
		c := &Compactor{windowFn: func() int { return tc.window }}
		if got := c.compressionInputBudget(); got != tc.want {
			t.Errorf("window %d: budget = %d, want %d", tc.window, got, tc.want)
		}
	}
}

func TestSplitCompressionInputPreservesAllUTF8Content(t *testing.T) {
	input := strings.Repeat("事实 #12：观察到服务响应。\n", 100)
	parts := splitCompressionInput(input, 128)
	if len(parts) < 2 || strings.Join(parts, "") != input {
		t.Fatal("split must preserve the complete source in order")
	}
	for _, part := range parts {
		if len(part) > 128 || !utf8.ValidString(part) {
			t.Fatalf("invalid part: %q", part)
		}
	}
}

func TestGroupCompressionSummariesShrinksOrFails(t *testing.T) {
	groups, err := groupCompressionSummaries([]string{"one", "two", "three"}, 40)
	if err != nil || len(groups) >= 3 {
		t.Fatalf("expected fewer groups: %v, %v", groups, err)
	}
	if _, err := groupCompressionSummaries([]string{strings.Repeat("x", 100)}, 40); err == nil {
		t.Fatal("oversized summary must fail instead of being truncated")
	}
}
