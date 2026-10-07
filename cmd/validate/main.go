package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	providercatalog "github.com/leamout/ai-providers/catalog"
	agentcontract "github.com/leamout/contracts/agent"
	"github.com/leamout/contracts/ai"
)

func main() {
	if err := validateRepository("."); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("agent manifests, tools, routing, and provider configurations are compatible")
}

func validateRepository(root string) error {
	providers := make(map[string]ai.Provider)
	for _, provider := range providercatalog.Providers() {
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
		if entry.IsDir() || entry.Name() != "agent.json" {
			return nil
		}
		manifest, err := readManifest(path)
		if err != nil {
			return err
		}
		return validateProviderBindings(path, manifest, providers)
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
	if err := agentcontract.ValidateTools(tools); err != nil {
		return fmt.Errorf("tools.json: %w", err)
	}
	toolNames := make(map[string]struct{}, len(tools))
	for _, tool := range tools {
		toolNames[tool.Name] = struct{}{}
	}

	aliases := make(map[string]struct{})
	agentFiles, err := filepath.Glob(filepath.Join(root, "agents", "*.json"))
	if err != nil {
		return err
	}
	if len(agentFiles) == 0 {
		return fmt.Errorf("template has no agent manifests")
	}
	for _, path := range agentFiles {
		manifest, err := readManifest(path)
		if err != nil {
			return err
		}
		if err := validateProviderBindings(path, manifest, providers); err != nil {
			return err
		}
		for _, tool := range manifest.Tools {
			if _, ok := toolNames[tool]; !ok {
				return fmt.Errorf("%s references unknown tool %q", path, tool)
			}
		}
		aliases[strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))] = struct{}{}
	}

	routingPath := filepath.Join(root, "routing.json")
	value, err := os.ReadFile(routingPath)
	if err != nil {
		return err
	}
	var routing agentcontract.RoutingManifest
	if err := json.Unmarshal(value, &routing); err != nil {
		return fmt.Errorf("%s: %w", routingPath, err)
	}
	if err := routing.Validate(); err != nil {
		return fmt.Errorf("%s: %w", routingPath, err)
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
			return fmt.Errorf("%s: official provider catalog does not contain %s", path, key)
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
