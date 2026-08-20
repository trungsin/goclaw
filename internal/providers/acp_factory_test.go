package providers

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewACPProviderFromRecord_RejectsRelativeBinary(t *testing.T) {
	_, err := NewACPProviderFromRecord("evil", "bash", nil, nil)
	if err == nil {
		t.Fatal("expected error for non-allowlisted relative binary")
	}
}

func TestNewACPProviderFromRecord_GrokAbsPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("posix stub")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "grok")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	prov, err := NewACPProviderFromRecord("grok-build", path, []byte(`{"perm_mode":"deny-all"}`), nil)
	if err != nil {
		t.Fatalf("NewACPProviderFromRecord: %v", err)
	}
	if prov.Name() != "grok-build" {
		t.Fatalf("Name()=%q", prov.Name())
	}
	if prov.DefaultModel() != "grok-build" {
		t.Fatalf("DefaultModel()=%q", prov.DefaultModel())
	}
	prov.Close()
}
