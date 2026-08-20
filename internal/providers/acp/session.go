package acp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"
)

func resolveACPWorkDir(workDir string) string {
	if workDir == "" {
		workDir = "."
	}
	abs, err := filepath.Abs(workDir)
	if err != nil {
		return workDir
	}
	return abs
}

// Initialize sends the ACP initialize request to establish capabilities.
func (p *ACPProcess) Initialize(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	req := InitializeRequest{
		ProtocolVersion: 1,
		ClientInfo:      ClientInfo{Name: "GoClaw", Version: "1.0"},
		Capabilities: ClientCaps{
			Fs:       &FsCaps{ReadTextFile: true, WriteTextFile: true},
			Terminal: true,
		},
	}
	var resp InitializeResponse
	if err := p.conn.Call(ctx, "initialize", req, &resp); err != nil {
		return fmt.Errorf("acp initialize: %w", err)
	}
	p.agentCaps = resp.Capabilities
	p.authMethods = resp.AuthMethods
	slog.Info("acp: initialized", "agent", resp.AgentInfo.Name, "version", resp.AgentInfo.Version, "loadSession", resp.Capabilities.LoadSession, "authMethods", len(resp.AuthMethods))
	return nil
}

// Authenticate completes the ACP auth handshake when the agent advertised methods.
func (p *ACPProcess) Authenticate(ctx context.Context, methodID string) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	req := AuthenticateRequest{
		MethodID: methodID,
		Meta:     map[string]any{"headless": true},
	}
	var resp json.RawMessage
	if err := p.conn.Call(ctx, "authenticate", req, &resp); err != nil {
		return fmt.Errorf("acp authenticate: %w", err)
	}
	slog.Info("acp: authenticated", "method", methodID)
	return nil
}

// NewSession creates a new ACP session and returns its session ID.
func (p *ACPProcess) NewSession(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cwd := resolveACPWorkDir(p.workDir)

	req := NewSessionRequest{
		Cwd:        cwd,
		McpServers: []string{},
	}
	var resp NewSessionResponse
	if err := p.conn.Call(ctx, "session/new", req, &resp); err != nil {
		return "", fmt.Errorf("acp session/new: %w", err)
	}
	slog.Info("acp: session/new", "sid", resp.SessionID, "cwd", cwd)
	return resp.SessionID, nil
}

// LoadSession restores a previous ACP session by ID (used after process restart).
// Returns the session ID to use going forward (may equal the requested ID).
// Only call if AgentCaps().LoadSession is true.
func (p *ACPProcess) LoadSession(ctx context.Context, sessionID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req := LoadSessionRequest{SessionID: sessionID, Cwd: resolveACPWorkDir(p.workDir)}
	var resp LoadSessionResponse
	if err := p.conn.Call(ctx, "session/load", req, &resp); err != nil {
		return "", fmt.Errorf("acp session/load: %w", err)
	}
	slog.Info("acp: session/load", "sid", resp.SessionID)
	return resp.SessionID, nil
}

// Prompt sends user content to sessionID and blocks until the agent completes,
// invoking onUpdate for each session/update notification received.
func (p *ACPProcess) Prompt(ctx context.Context, sessionID string, content []ContentBlock, onUpdate func(SessionUpdate)) (*PromptResponse, error) {
	p.inUse.Add(1)
	defer p.inUse.Add(-1)

	p.mu.Lock()
	p.lastActive = time.Now()
	p.mu.Unlock()

	p.registerUpdateFn(sessionID, onUpdate)
	defer p.unregisterUpdateFn(sessionID)

	goclawSession := goclawSessionFromCtx(ctx)
	slog.Info("acp: session/prompt", "session", goclawSession, "sid", sessionID)
	req := PromptRequest{
		SessionID: sessionID,
		Prompt:    content,
	}

	var resp PromptResponse
	if err := p.conn.Call(ctx, "session/prompt", req, &resp); err != nil {
		return nil, fmt.Errorf("acp session/prompt: %w", err)
	}

	p.mu.Lock()
	p.lastActive = time.Now()
	p.mu.Unlock()

	slog.Info("acp: session/prompt completed", "session", goclawSession, "sid", sessionID, "stopReason", resp.StopReason)
	return &resp, nil
}

// Cancel sends a session/cancel notification for the given session.
func (p *ACPProcess) Cancel(sessionID string) error {
	return p.conn.Notify("session/cancel", CancelNotification{
		SessionID: sessionID,
	})
}
