package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/leamout/ai-providers/assemblyai"
	"github.com/leamout/ai-providers/cartesia"
	"github.com/leamout/ai-providers/deepgram"
	"github.com/leamout/ai-providers/elevenlabs"
	"github.com/leamout/ai-providers/gemini"
	"github.com/leamout/ai-providers/groq"
	"github.com/leamout/ai-providers/openai"
	agentcontract "github.com/leamout/contracts/agent"
	"github.com/leamout/contracts/ai"
)

// Match the portable package limits enforced by the Leamout runtime.
const maxPortableAgents = 32
const maxPortableTools = 128

var agentAliasPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

func main() {
	if err := validateRepository("."); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("agent manifests, tools, routing, and provider configurations are compatible")
}

func validateRepository(root string) error {
	providers := make(map[string]ai.Provider)
	for _, provider := range validationProviders() {
		descriptor := provider.Descriptor()
		providers[providerKey(descriptor.Kind, descriptor.ID)] = provider
	}

	if err := validateExamples(filepath.Join(root, "examples"), providers); err != nil {
		return err
	}
	return validateTemplates(filepath.Join(root, "templates"), providers)
}

func validateExamples(root string, providers map[string]ai.Provider) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		switch entry.Name() {
		case "agent.json":
			manifest, err := readManifest(path)
			if err != nil {
				return err
			}
			return validateProviderBindings(path, manifest, providers)
		case "routing.json":
			_, err := readRouting(path)
			return err
		default:
			return nil
		}
	})
}

func validateTemplates(root string, providers map[string]ai.Provider) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		templateRoot := filepath.Join(root, entry.Name())
		if err := validateTemplate(templateRoot, providers); err != nil {
			return fmt.Errorf("%s: %w", entry.Name(), err)
		}
	}
	return nil
}

func validateTemplate(root string, providers map[string]ai.Provider) error {
	tools, err := readTools(filepath.Join(root, "tools.json"))
	if err != nil {
		return err
	}
	if len(tools) > maxPortableTools {
		return fmt.Errorf("package cannot contain more than %d tools", maxPortableTools)
	}
	if err := agentcontract.ValidateTools(tools); err != nil {
		return fmt.Errorf("tools.json: %w", err)
	}

	toolNames := make(map[string]struct{}, len(tools))
	for _, tool := range tools {
		toolNames[strings.TrimSpace(tool.Name)] = struct{}{}
	}

	entries, err := os.ReadDir(filepath.Join(root, "agents"))
	if err != nil {
		return fmt.Errorf("read agents directory: %w", err)
	}
	if len(entries) == 0 {
		return fmt.Errorf("package must contain at least one agent")
	}
	if len(entries) > maxPortableAgents {
		return fmt.Errorf("package cannot contain more than %d agents", maxPortableAgents)
	}

	aliases := make(map[string]struct{}, len(entries))
	names := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			return fmt.Errorf("unexpected agents entry %q", entry.Name())
		}
		alias := strings.TrimSuffix(entry.Name(), ".json")
		if !agentAliasPattern.MatchString(alias) {
			return fmt.Errorf("invalid agent alias %q", alias)
		}

		path := filepath.Join(root, "agents", entry.Name())
		manifest, err := readManifest(path)
		if err != nil {
			return err
		}
		if err := validateProviderBindings(path, manifest, providers); err != nil {
			return err
		}
		for _, tool := range manifest.Tools {
			if _, ok := toolNames[strings.TrimSpace(tool)]; !ok {
				return fmt.Errorf("%s references unknown tool %q", path, tool)
			}
		}

		name := strings.ToLower(strings.TrimSpace(manifest.Name))
		if _, exists := names[name]; exists {
			return fmt.Errorf("duplicate agent name %q", manifest.Name)
		}
		names[name] = struct{}{}
		aliases[alias] = struct{}{}
	}

	routingPath := filepath.Join(root, "routing.json")
	routing, err := readRouting(routingPath)
	if err != nil {
		return err
	}
	if err := routing.ValidateAgentReferences(aliases); err != nil {
		return fmt.Errorf("%s: %w", routingPath, err)
	}
	return nil
}

func readManifest(path string) (agentcontract.Manifest, error) {
	value, err := os.ReadFile(path)
	if err != nil {
		return agentcontract.Manifest{}, err
	}
	var manifest agentcontract.Manifest
	if err := json.Unmarshal(value, &manifest); err != nil {
		return agentcontract.Manifest{}, fmt.Errorf("%s: %w", path, err)
	}
	if err := manifest.Validate(); err != nil {
		return agentcontract.Manifest{}, fmt.Errorf("%s: %w", path, err)
	}
	return manifest, nil
}

func readRouting(path string) (agentcontract.RoutingManifest, error) {
	value, err := os.ReadFile(path)
	if err != nil {
		return agentcontract.RoutingManifest{}, err
	}
	var routing agentcontract.RoutingManifest
	if err := json.Unmarshal(value, &routing); err != nil {
		return agentcontract.RoutingManifest{}, fmt.Errorf("%s: %w", path, err)
	}
	if err := routing.Validate(); err != nil {
		return agentcontract.RoutingManifest{}, fmt.Errorf("%s: %w", path, err)
	}
	return routing, nil
}

func readTools(path string) ([]agentcontract.ToolDefinition, error) {
	value, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tools []agentcontract.ToolDefinition
	if err := json.Unmarshal(value, &tools); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return tools, nil
}

func validateProviderBindings(path string, manifest agentcontract.Manifest, providers map[string]ai.Provider) error {
	for _, binding := range manifest.Providers {
		key := providerKey(binding.Role, binding.Provider)
		provider, ok := providers[key]
		if !ok {
			return fmt.Errorf("%s: validation adapter set does not contain %s", path, key)
		}
		if len(binding.Config) == 0 {
			continue
		}
		validator, ok := provider.(ai.ConfigValidator)
		if !ok {
			continue
		}
		if err := validator.ValidateConfig(binding.Config); err != nil {
			return fmt.Errorf("%s: invalid %s config: %w", path, key, err)
		}
	}
	return nil
}

func providerKey(kind ai.Kind, id string) string {
	return string(kind) + ":" + strings.TrimSpace(id)
}

func validationProviders() []ai.Provider {
	return []ai.Provider{
		deepgram.Provider{},
		assemblyai.Provider{},
		groq.Provider{},
		openai.Provider{},
		openai.RealtimeProvider{},
		cartesia.Provider{},
		elevenlabs.Provider{},
		gemini.RealtimeProvider{},
	}
}
