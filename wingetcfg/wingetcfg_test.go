package wingetcfg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWingetCfg(t *testing.T) {
	cfg := NewWingetCfg()
	if cfg == nil {
		t.Fatal("NewWingetCfg returned nil")
	}
	if cfg.Properties.ConfigurationVersion != WinGetConfigurationVersion {
		t.Errorf("expected ConfigurationVersion %q, got %q", WinGetConfigurationVersion, cfg.Properties.ConfigurationVersion)
	}
	if cfg.Properties.Resources != nil {
		t.Error("expected Resources to be nil initially")
	}
	if cfg.Properties.Assertions != nil {
		t.Error("expected Assertions to be nil initially")
	}
}

func TestAddResource(t *testing.T) {
	cfg := NewWingetCfg()
	r := &WinGetResource{Resource: "test/resource"}
	cfg.AddResource(r)

	if len(cfg.Properties.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(cfg.Properties.Resources))
	}
	if cfg.Properties.Resources[0].Resource != "test/resource" {
		t.Errorf("expected resource %q, got %q", "test/resource", cfg.Properties.Resources[0].Resource)
	}
}

func TestAddMultipleResources(t *testing.T) {
	cfg := NewWingetCfg()
	cfg.AddResource(&WinGetResource{Resource: "res1"})
	cfg.AddResource(&WinGetResource{Resource: "res2"})
	cfg.AddResource(&WinGetResource{Resource: "res3"})

	if len(cfg.Properties.Resources) != 3 {
		t.Fatalf("expected 3 resources, got %d", len(cfg.Properties.Resources))
	}
}

func TestAddAssertion(t *testing.T) {
	cfg := NewWingetCfg()
	r := &WinGetResource{Resource: "test/assertion"}
	cfg.AddAssertion(r)

	if len(cfg.Properties.Assertions) != 1 {
		t.Fatalf("expected 1 assertion, got %d", len(cfg.Properties.Assertions))
	}
	if cfg.Properties.Assertions[0].Resource != "test/assertion" {
		t.Errorf("expected assertion %q, got %q", "test/assertion", cfg.Properties.Assertions[0].Resource)
	}
}

func TestAddMultipleAssertions(t *testing.T) {
	cfg := NewWingetCfg()
	cfg.AddAssertion(&WinGetResource{Resource: "assert1"})
	cfg.AddAssertion(&WinGetResource{Resource: "assert2"})

	if len(cfg.Properties.Assertions) != 2 {
		t.Fatalf("expected 2 assertions, got %d", len(cfg.Properties.Assertions))
	}
}

func TestWriteConfigFile(t *testing.T) {
	cfg := NewWingetCfg()
	r, err := InstallPackage("pkg1", "Install Firefox", "Mozilla.Firefox", "winget", "", true)
	if err != nil {
		t.Fatalf("unexpected error creating package: %v", err)
	}
	cfg.AddResource(r)

	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "test_wingetcfg_output.yaml")
	defer os.Remove(filePath)

	err = cfg.WriteConfigFile(filePath)
	if err != nil {
		t.Fatalf("WriteConfigFile returned error: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	content := string(data)

	if !strings.HasPrefix(content, DSCSchema) {
		t.Error("output file does not start with DSC schema header")
	}
	if !strings.Contains(content, "configurationVersion: 0.2.0") {
		t.Error("output file does not contain configurationVersion")
	}
	if !strings.Contains(content, "Mozilla.Firefox") {
		t.Error("output file does not contain package ID")
	}
}

func TestWriteConfigFileEmpty(t *testing.T) {
	cfg := NewWingetCfg()

	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "test_wingetcfg_empty.yaml")
	defer os.Remove(filePath)

	err := cfg.WriteConfigFile(filePath)
	if err != nil {
		t.Fatalf("WriteConfigFile returned error: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	content := string(data)

	if !strings.HasPrefix(content, DSCSchema) {
		t.Error("output file does not start with DSC schema header")
	}
	if !strings.Contains(content, "configurationVersion: 0.2.0") {
		t.Error("output file does not contain configurationVersion")
	}
}

func TestWriteConfigFileInvalidPath(t *testing.T) {
	cfg := NewWingetCfg()
	err := cfg.WriteConfigFile("/nonexistent/directory/file.yaml")
	if err == nil {
		t.Error("expected error for invalid file path, got nil")
	}
}

func TestSetEnsureValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Present", input: EnsurePresent, expected: EnsurePresent},
		{name: "Absent", input: EnsureAbsent, expected: EnsureAbsent},
		{name: "Empty string defaults to Present", input: "", expected: EnsurePresent},
		{name: "Invalid value defaults to Present", input: "Invalid", expected: EnsurePresent},
		{name: "Lowercase present defaults to Present", input: "present", expected: EnsurePresent},
		{name: "Lowercase absent defaults to Present", input: "absent", expected: EnsurePresent},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := SetEnsureValue(tc.input)
			if result != tc.expected {
				t.Errorf("SetEnsureValue(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if DSCSchema != "# yaml-language-server: $schema=https://aka.ms/configuration-dsc-schema/0.2" {
		t.Errorf("unexpected DSCSchema value: %q", DSCSchema)
	}
	if WinGetConfigurationVersion != "0.2.0" {
		t.Errorf("unexpected WinGetConfigurationVersion value: %q", WinGetConfigurationVersion)
	}
	if EnsurePresent != "Present" {
		t.Errorf("unexpected EnsurePresent value: %q", EnsurePresent)
	}
	if EnsureAbsent != "Absent" {
		t.Errorf("unexpected EnsureAbsent value: %q", EnsureAbsent)
	}
}
