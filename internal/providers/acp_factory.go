package providers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/nextlevelbuilder/goclaw/internal/config"
	"github.com/nextlevelbuilder/goclaw/internal/providers/acp"
)

// DefaultACPWorkDir is the workspace used when a provider record omits work_dir.
func DefaultACPWorkDir() string {
	return filepath.Join(config.ResolvedDataDirFromEnv(), "acp-workspaces")
}

type acpRecordSettings struct {
	Args     []string `json:"args"`
	IdleTTL  string   `json:"idle_ttl"`
	PermMode string   `json:"perm_mode"`
	WorkDir  string   `json:"work_dir"`
}

// NewACPProviderFromRecord builds an ACP provider from a DB/config record.
// binary is stored in api_base. Returns an error if the binary is rejected or missing.
func NewACPProviderFromRecord(name, binary string, settings json.RawMessage, deny []*regexp.Regexp) (*ACPProvider, error) {
	if binary == "" {
		return nil, fmt.Errorf("acp: no binary specified")
	}
	if !acp.AllowedBinary(binary) {
		return nil, fmt.Errorf("acp: invalid binary path")
	}
	if _, err := exec.LookPath(binary); err != nil {
		return nil, fmt.Errorf("acp: binary not found: %s", binary)
	}
	var s acpRecordSettings
	if len(settings) > 0 {
		if err := json.Unmarshal(settings, &s); err != nil {
			slog.Warn("acp: invalid settings JSON, using defaults", "name", name, "error", err)
		}
	}
	idleTTL := 5 * time.Minute
	if s.IdleTTL != "" {
		if d, err := time.ParseDuration(s.IdleTTL); err == nil {
			idleTTL = d
		}
	}
	workDir := s.WorkDir
	if workDir == "" {
		workDir = DefaultACPWorkDir()
	}
	opts := []ACPOption{WithACPName(name), WithACPModel(name)}
	if s.PermMode != "" {
		opts = append(opts, WithACPPermMode(s.PermMode))
	}
	return NewACPProvider(binary, acp.DefaultArgs(binary, s.Args), workDir, idleTTL, deny, opts...), nil
}
