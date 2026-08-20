package acp

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestGrokACPHandshake(t *testing.T) {
	if os.Getenv("ACP_GROK_E2E") == "" {
		t.Skip("set ACP_GROK_E2E=1 to run this test (requires grok binary + login)")
	}
	if _, err := exec.LookPath("grok"); err != nil {
		t.Skip("grok binary not on PATH")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	pp := NewProcessPool("grok", DefaultArgs("grok", nil), ".", 10*time.Minute)
	defer pp.Close()

	proc, err := pp.GetOrSpawn(ctx, "grok-e2e")
	if err != nil {
		t.Fatalf("Spawn failed: %v", err)
	}

	sid, err := proc.NewSession(ctx)
	if err != nil {
		t.Fatalf("NewSession failed: %v", err)
	}
	if sid == "" {
		t.Fatal("empty session id")
	}

	var collected strings.Builder
	_, err = proc.Prompt(ctx, sid, []ContentBlock{
		{Type: "text", Text: "Reply with ACP_GROK_OK and nothing else."},
	}, func(su SessionUpdate) {
		if su.Message == nil {
			return
		}
		for _, b := range su.Message.Content {
			if b.Type == "text" {
				collected.WriteString(b.Text)
			}
		}
	})
	if err != nil {
		t.Fatalf("Prompt failed: %v", err)
	}
	if collected.String() == "" {
		t.Error("no text collected from grok ACP session/update")
	}
}
