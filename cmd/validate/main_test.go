package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	agentcontract "github.com/leamout/contracts/agent"
	"github.com/leamout/contracts/ai"
)

type templateFixture struct {
	agents                  map[string]agentcontract.Manifest
	tools                   []agentcontract.ToolDefinition
	routing                 agentcontract.RoutingManifest
	unexpectedAgentFilename string
}

func validTemplateFixture() templateFixture {
	return templateFixture{
		agents: map[string]agentcontract.Manifest{
			"main": validTestManifest("Main Agent"),
		},
		routing: agentcontract.RoutingManifest{
			SchemaVersion: agentcontract.SchemaVersion,
			Name:          "Example routing",
			DefaultAgent:  "main",
		},
	}
}

func validTestManifest(name string) agentcontract.Manifest {
	return agentcontract.Manifest{
		SchemaVersion: agentcontract.SchemaVersion,
		Name:          name,
		Engine:        agentcontract.EngineRealtime,
		Instructions:  "Help the caller.",
		Providers: []agentcontract.ProviderBinding{
			{
				Role:     ai.KindRealtime,
				Provider: "openai",
				Config:   json.RawMessage(`{"model":"gpt-realtime-2.1"}`),
			},
		},
	}
}

func TestTemplateValidationMatchesRuntime(t *testing.T) {
	providers := make(map[string]ai.Provider)
	for _, provider := range validationProviders() {
		descriptor := provider.Descriptor()
		providers[providerKey(descriptor.Kind, descriptor.ID)] = provider
	}

	tests := []struct {
		name   string
		mutate func(*templateFixture)
		want   string
	}{
		{
			name: "valid portable package",
		},
		{
			name: "missing agents",
			mutate: func(f *templateFixture) {
				f.agents = nil
			},
			want: "at least one agent",
		},
		{
			name: "too many agents",
			mutate: func(f *templateFixture) {
				for i := 0; i < maxPortableAgents; i++ {
					name := fmt.Sprintf("extra-%02d", i)
					f.agents[name] = validTestManifest(name)
				}
			},
			want: "more than 32 agents",
		},
		{
			name: "too many tools",
			mutate: func(f *templateFixture) {
				for i := 0; i <= maxPortableTools; i++ {
					f.tools = append(f.tools, agentcontract.ToolDefinition{
						Type:        agentcontract.ToolTypeWebhook,
						Name:        fmt.Sprintf("tool_%d", i),
						Description: "Example tool",
						EndpointURL: "${TOOL_ENDPOINT}",
					})
				}
			},
			want: "more than 128 tools",
		},
		{
			name: "invalid routing alias",
			mutate: func(f *templateFixture) {
				f.agents["Uppercase"] = validTestManifest("Another Agent")
			},
			want: "invalid agent alias",
		},
		{
			name: "duplicate agent names ignoring case and whitespace",
			mutate: func(f *templateFixture) {
				f.agents["secondary"] = validTestManifest(" main agent ")
			},
			want: "duplicate agent name",
		},
		{
			name: "unexpected agent file",
			mutate: func(f *templateFixture) {
				f.unexpectedAgentFilename = "notes.txt"
			},
			want: "unexpected agents entry",
		},
		{
			name: "unknown tool reference",
			mutate: func(f *templateFixture) {
				agent := f.agents["main"]
				agent.Tools = []string{"unknown_tool"}
				f.agents["main"] = agent
			},
			want: "references unknown tool",
		},
		{
			name: "unknown provider",
			mutate: func(f *templateFixture) {
				agent := f.agents["main"]
				agent.Providers[0].Provider = "unregistered"
				f.agents["main"] = agent
			},
			want: "validation adapter set does not contain",
		},
		{
			name: "unapproved provider endpoint",
			mutate: func(f *templateFixture) {
				agent := f.agents["main"]
				agent.Providers[0].Config = json.RawMessage(`{"endpoint":"wss://untrusted.example/ws"}`)
				f.agents["main"] = agent
			},
			want: "invalid realtime:openai config",
		},
		{
			name: "unknown routing reference",
			mutate: func(f *templateFixture) {
				f.routing.DefaultAgent = "missing"
			},
			want: "unknown default agent",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := validTemplateFixture()
			if test.mutate != nil {
				test.mutate(&fixture)
			}
			root := t.TempDir()
			writeTemplateFixture(t, root, fixture)

			err := validateTemplate(root, providers)
			if test.want == "" {
				if err != nil {
					t.Fatalf("valid template rejected: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected %q validation error, got %v", test.want, err)
			}
		})
	}
}

func TestExampleRoutingIsValidated(t *testing.T) {
	root := t.TempDir()
	writeTestJSON(t, filepath.Join(root, "scheduled-routing", "routing.json"), agentcontract.RoutingManifest{
		SchemaVersion: agentcontract.SchemaVersion,
		Name:          "Invalid schedule",
		DefaultAgent:  "main",
		Routes: []agentcontract.Route{
			{Schedule: "missing", Agent: "main"},
		},
	})

	if err := validateExamples(root, nil); err == nil ||
		!strings.Contains(err.Error(), "unknown schedule") {
		t.Fatalf("expected invalid example routing to be rejected, got %v", err)
	}
}

func writeTemplateFixture(t *testing.T, root string, fixture templateFixture) {
	t.Helper()
	agentsDir := filepath.Join(root, "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for alias, manifest := range fixture.agents {
		writeTestJSON(t, filepath.Join(agentsDir, alias+".json"), manifest)
	}
	if fixture.unexpectedAgentFilename != "" {
		if err := os.WriteFile(
			filepath.Join(agentsDir, fixture.unexpectedAgentFilename),
			[]byte("unexpected"),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}
	writeTestJSON(t, filepath.Join(root, "tools.json"), fixture.tools)
	writeTestJSON(t, filepath.Join(root, "routing.json"), fixture.routing)
}

func writeTestJSON(t *testing.T, path string, value any) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatal(err)
	}
}
