package wingetcfg

import "testing"

func TestNewWinGetPackageResource(t *testing.T) {
	tests := []struct {
		name            string
		id              string
		description     string
		packageID       string
		source          string
		version         string
		useLatest       bool
		ensure          bool
		expectErr       bool
		expectSource    string
		expectEnsure    string
		expectVersion   string
		expectUseLatest bool
	}{
		{
			name:            "Valid install with defaults",
			id:              "pkg1",
			description:     "Install package",
			packageID:       "Mozilla.Firefox",
			source:          "",
			version:         "",
			useLatest:       true,
			ensure:          true,
			expectErr:       false,
			expectSource:    "winget",
			expectEnsure:    "Present",
			expectUseLatest: true,
		},
		{
			name:            "Valid install with custom source",
			id:              "pkg2",
			description:     "Install from store",
			packageID:       "Mozilla.Firefox",
			source:          "msstore",
			version:         "",
			useLatest:       true,
			ensure:          true,
			expectErr:       false,
			expectSource:    "msstore",
			expectEnsure:    "Present",
			expectUseLatest: true,
		},
		{
			name:            "Valid install with specific version",
			id:              "pkg3",
			description:     "Install specific version",
			packageID:       "Mozilla.Firefox",
			source:          "winget",
			version:         "100.0",
			useLatest:       false,
			ensure:          true,
			expectErr:       false,
			expectSource:    "winget",
			expectEnsure:    "Present",
			expectVersion:   "100.0",
			expectUseLatest: false,
		},
		{
			name:            "Version ignored when useLatest is true",
			id:              "pkg4",
			description:     "Use latest overrides version",
			packageID:       "Mozilla.Firefox",
			source:          "",
			version:         "100.0",
			useLatest:       true,
			ensure:          true,
			expectErr:       false,
			expectSource:    "winget",
			expectEnsure:    "Present",
			expectUseLatest: true,
		},
		{
			name:         "Uninstall package",
			id:           "pkg5",
			description:  "Uninstall package",
			packageID:    "Mozilla.Firefox",
			source:       "",
			version:      "",
			useLatest:    false,
			ensure:       false,
			expectErr:    false,
			expectSource: "winget",
			expectEnsure: "Absent",
		},
		{
			name:        "Empty packageID returns error",
			id:          "pkg6",
			description: "Should fail",
			packageID:   "",
			source:      "",
			version:     "",
			useLatest:   true,
			ensure:      true,
			expectErr:   true,
		},
		{
			name:            "Empty ID is allowed",
			id:              "",
			description:     "No ID",
			packageID:       "Mozilla.Firefox",
			source:          "",
			version:         "",
			useLatest:       true,
			ensure:          true,
			expectErr:       false,
			expectSource:    "winget",
			expectEnsure:    "Present",
			expectUseLatest: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewWinGetPackageResource(tc.id, tc.description, tc.packageID, tc.source, tc.version, tc.useLatest, tc.ensure)
			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if r.Resource != WinGetPackageResource {
				t.Errorf("expected resource %q, got %q", WinGetPackageResource, r.Resource)
			}
			if tc.id != "" && r.ID != tc.id {
				t.Errorf("expected ID %q, got %q", tc.id, r.ID)
			}
			if tc.id == "" && r.ID != "" {
				t.Errorf("expected empty ID, got %q", r.ID)
			}
			if r.Directives.Description != tc.description {
				t.Errorf("expected description %q, got %q", tc.description, r.Directives.Description)
			}
			if !r.Directives.AllowPreRelease {
				t.Error("expected AllowPreRelease to be true")
			}
			if r.Settings["source"] != tc.expectSource {
				t.Errorf("expected source %q, got %q", tc.expectSource, r.Settings["source"])
			}
			if r.Settings["Ensure"] != tc.expectEnsure {
				t.Errorf("expected Ensure %q, got %q", tc.expectEnsure, r.Settings["Ensure"])
			}
			if r.Settings["id"] != tc.packageID {
				t.Errorf("expected id %q, got %q", tc.packageID, r.Settings["id"])
			}
			if tc.expectVersion != "" {
				if r.Settings["version"] != tc.expectVersion {
					t.Errorf("expected version %q, got %v", tc.expectVersion, r.Settings["version"])
				}
			}
			if r.Settings["uselatest"] != tc.expectUseLatest {
				t.Errorf("expected uselatest %v, got %v", tc.expectUseLatest, r.Settings["uselatest"])
			}
		})
	}
}

func TestInstallPackage(t *testing.T) {
	r, err := InstallPackage("id1", "Install Firefox", "Mozilla.Firefox", "winget", "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != "Present" {
		t.Errorf("expected Ensure Present, got %v", r.Settings["Ensure"])
	}
}

func TestInstallPackageEmptyID(t *testing.T) {
	_, err := InstallPackage("id1", "Fail", "", "winget", "", true)
	if err == nil {
		t.Fatal("expected error for empty packageID, got nil")
	}
}

func TestUninstallPackage(t *testing.T) {
	r, err := UninstallPackage("id1", "Uninstall Firefox", "Mozilla.Firefox", "winget", "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != "Absent" {
		t.Errorf("expected Ensure Absent, got %v", r.Settings["Ensure"])
	}
}

func TestUninstallPackageEmptyID(t *testing.T) {
	_, err := UninstallPackage("id1", "Fail", "", "winget", "", false)
	if err == nil {
		t.Fatal("expected error for empty packageID, got nil")
	}
}
