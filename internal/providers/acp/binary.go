package acp

import (
	"path/filepath"
	"strings"
)

func agentBaseName(binary string) string {
	return strings.TrimSuffix(strings.ToLower(filepath.Base(binary)), ".exe")
}

func knownAgentName(name string) bool {
	switch name {
	case "claude", "codex", "gemini", "grok":
		return true
	default:
		return false
	}
}

// AllowedBinary reports whether an ACP agent binary is safe to spawn.
// Relative names must be an allowlisted command. Absolute paths are accepted
// only when the basename is also allowlisted (blocks /bin/bash).
func AllowedBinary(binary string) bool {
	if binary == "" || strings.Contains(binary, "\x00") || strings.Contains(binary, "://") {
		return false
	}
	if filepath.IsAbs(binary) {
		return knownAgentName(agentBaseName(binary))
	}
	return knownAgentName(strings.TrimSuffix(strings.ToLower(binary), ".exe"))
}

// DefaultArgs returns spawn args when the caller left them empty.
// Grok Build speaks ACP only via `grok agent stdio`.
func DefaultArgs(binary string, args []string) []string {
	if len(args) > 0 {
		return args
	}
	if agentBaseName(binary) == "grok" {
		return []string{"agent", "stdio"}
	}
	return args
}
