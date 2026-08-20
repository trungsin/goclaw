package acp

import "fmt"

// AuthMethod is an ACP authenticate method advertised in initialize.
type AuthMethod struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// AuthenticateRequest is the client → agent authenticate call.
type AuthenticateRequest struct {
	MethodID string         `json:"methodId"`
	Meta     map[string]any `json:"_meta,omitempty"`
}

func resolveXAIKey(getenv func(string) string) string {
	if getenv == nil {
		return ""
	}
	if v := getenv("XAI_API_KEY"); v != "" {
		return v
	}
	return getenv("GOCLAW_XAI_API_KEY")
}

// AuthMethodCandidates returns headless auth methods in retry order.
// Empty methods means the agent needs no authenticate call.
func AuthMethodCandidates(methods []AuthMethod, getenv func(string) string) ([]string, error) {
	if len(methods) == 0 {
		return nil, nil
	}
	has := make(map[string]struct{}, len(methods))
	for _, m := range methods {
		if m.ID == "" {
			continue
		}
		has[m.ID] = struct{}{}
	}
	var out []string
	if _, ok := has["cached_token"]; ok {
		out = append(out, "cached_token")
	}
	if _, ok := has["xai.api_key"]; ok && resolveXAIKey(getenv) != "" {
		out = append(out, "xai.api_key")
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("acp: no headless auth method (run grok login or set XAI_API_KEY)")
	}
	return out, nil
}

// PickAuthMethod returns the first headless-capable ACP auth method.
func PickAuthMethod(methods []AuthMethod, getenv func(string) string) (string, error) {
	cands, err := AuthMethodCandidates(methods, getenv)
	if err != nil || len(cands) == 0 {
		return "", err
	}
	return cands[0], nil
}
