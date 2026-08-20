package acp

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestAllowedBinary(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", false},
		{"claude", true},
		{"codex", true},
		{"gemini", true},
		{"grok", true},
		{"bash", false},
		{"grok.exe", true},
		{"/usr/bin/grok", true},
		{"/bin/bash", false},
		{"https://example.com/grok", false},
		{"relative/grok", false},
		{filepath.Join(string(filepath.Separator), "opt", "grok"), true},
	}
	for _, tc := range cases {
		if got := AllowedBinary(tc.in); got != tc.want {
			t.Errorf("AllowedBinary(%q)=%v want %v", tc.in, got, tc.want)
		}
	}
}

func TestDefaultArgs(t *testing.T) {
	if got := DefaultArgs("grok", nil); !reflect.DeepEqual(got, []string{"agent", "stdio"}) {
		t.Fatalf("grok empty args: %v", got)
	}
	custom := []string{"--always-approve", "agent", "stdio"}
	if got := DefaultArgs("grok", custom); !reflect.DeepEqual(got, custom) {
		t.Fatalf("grok custom args overwritten: %v", got)
	}
	abs := filepath.Join(string(filepath.Separator), "Users", "bin", "grok")
	if got := DefaultArgs(abs, nil); !reflect.DeepEqual(got, []string{"agent", "stdio"}) {
		t.Fatalf("abs grok args: %v", got)
	}
	if got := DefaultArgs("claude", nil); got != nil {
		t.Fatalf("claude empty args should stay nil, got %v", got)
	}
}
