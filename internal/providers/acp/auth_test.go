package acp

import "testing"

func TestPickAuthMethod(t *testing.T) {
	cached := []AuthMethod{{ID: "cached_token"}, {ID: "xai.api_key"}, {ID: "grok.com"}}
	apiKey := []AuthMethod{{ID: "xai.api_key"}, {ID: "grok.com"}}
	browserOnly := []AuthMethod{{ID: "grok.com"}}
	other := []AuthMethod{{ID: "token"}}
	emptyID := []AuthMethod{{ID: ""}, {ID: "grok.com"}}

	env := func(key string) string {
		if key == "XAI_API_KEY" {
			return "xai-test"
		}
		return ""
	}
	goclawEnv := func(key string) string {
		if key == "GOCLAW_XAI_API_KEY" {
			return "xai-from-goclaw"
		}
		return ""
	}
	emptyEnv := func(string) string { return "" }

	t.Run("no methods", func(t *testing.T) {
		got, err := PickAuthMethod(nil, env)
		if err != nil || got != "" {
			t.Fatalf("got %q err %v, want empty", got, err)
		}
	})
	t.Run("cached_token preferred", func(t *testing.T) {
		got, err := PickAuthMethod(cached, emptyEnv)
		if err != nil || got != "cached_token" {
			t.Fatalf("got %q err %v", got, err)
		}
	})
	t.Run("xai.api_key when env set", func(t *testing.T) {
		got, err := PickAuthMethod(apiKey, env)
		if err != nil || got != "xai.api_key" {
			t.Fatalf("got %q err %v", got, err)
		}
	})
	t.Run("GOCLAW_XAI_API_KEY accepted", func(t *testing.T) {
		got, err := PickAuthMethod(apiKey, goclawEnv)
		if err != nil || got != "xai.api_key" {
			t.Fatalf("got %q err %v", got, err)
		}
	})
	t.Run("xai.api_key skipped without env", func(t *testing.T) {
		if _, err := PickAuthMethod(apiKey, emptyEnv); err == nil {
			t.Fatal("expected error for browser-only leftover")
		}
	})
	t.Run("browser only errors", func(t *testing.T) {
		if _, err := PickAuthMethod(browserOnly, emptyEnv); err == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("unknown method rejected", func(t *testing.T) {
		if _, err := PickAuthMethod(other, emptyEnv); err == nil {
			t.Fatal("expected error for unknown method")
		}
	})
	t.Run("empty method id ignored", func(t *testing.T) {
		if _, err := PickAuthMethod(emptyID, emptyEnv); err == nil {
			t.Fatal("expected error when only empty/browser methods remain")
		}
	})
}

func TestAuthMethodCandidatesFallbackOrder(t *testing.T) {
	methods := []AuthMethod{{ID: "cached_token"}, {ID: "xai.api_key"}}
	env := func(key string) string {
		if key == "XAI_API_KEY" {
			return "xai-test"
		}
		return ""
	}
	got, err := AuthMethodCandidates(methods, env)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "cached_token" || got[1] != "xai.api_key" {
		t.Fatalf("candidates = %v", got)
	}
}
