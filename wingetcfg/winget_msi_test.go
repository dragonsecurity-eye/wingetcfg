package wingetcfg

import "testing"

func TestNewMSIPackageResource(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		description   string
		productID     string
		path          string
		arguments     string
		logPath       string
		fileHash      string
		hashAlgorithm string
		ensure        bool
		expectErr     bool
		expectEnsure  string
		expectArgs    bool
		expectLogPath bool
	}{
		{
			name:          "Valid install with all fields",
			id:            "msi1",
			description:   "Install MSI",
			productID:     "{12345}",
			path:          "C:\\installer.msi",
			arguments:     "/quiet",
			logPath:       "C:\\log.txt",
			fileHash:      "abc123",
			hashAlgorithm: "SHA256",
			ensure:        true,
			expectErr:     false,
			expectEnsure:  "Present",
			expectArgs:    true,
			expectLogPath: true,
		},
		{
			name:          "Valid install minimal fields",
			id:            "",
			description:   "Install MSI",
			productID:     "",
			path:          "C:\\installer.msi",
			arguments:     "",
			logPath:       "",
			fileHash:      "",
			hashAlgorithm: "",
			ensure:        true,
			expectErr:     false,
			expectEnsure:  "Present",
			expectArgs:    false,
			expectLogPath: false,
		},
		{
			name:          "Uninstall MSI",
			id:            "msi2",
			description:   "Uninstall MSI",
			productID:     "{12345}",
			path:          "C:\\installer.msi",
			arguments:     "",
			logPath:       "",
			fileHash:      "",
			hashAlgorithm: "",
			ensure:        false,
			expectErr:     false,
			expectEnsure:  "Absent",
		},
		{
			name:          "Empty path returns error",
			id:            "msi3",
			description:   "Should fail",
			productID:     "{12345}",
			path:          "",
			arguments:     "",
			logPath:       "",
			fileHash:      "",
			hashAlgorithm: "",
			ensure:        true,
			expectErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewMSIPackageResource(tc.id, tc.description, tc.productID, tc.path, tc.arguments, tc.logPath, tc.fileHash, tc.hashAlgorithm, tc.ensure)
			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if r.Resource != WinGetMSIPackageResource {
				t.Errorf("expected resource %q, got %q", WinGetMSIPackageResource, r.Resource)
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
			if r.Settings["Ensure"] != tc.expectEnsure {
				t.Errorf("expected Ensure %q, got %v", tc.expectEnsure, r.Settings["Ensure"])
			}
			if r.Settings["Path"] != tc.path {
				t.Errorf("expected Path %q, got %v", tc.path, r.Settings["Path"])
			}

			_, hasArgs := r.Settings["Arguments"]
			if tc.expectArgs && !hasArgs {
				t.Error("expected Arguments to be set")
			}
			if !tc.expectArgs && hasArgs {
				t.Error("expected Arguments to not be set")
			}

			_, hasLogPath := r.Settings["LogPath"]
			if tc.expectLogPath && !hasLogPath {
				t.Error("expected LogPath to be set")
			}
			if !tc.expectLogPath && hasLogPath {
				t.Error("expected LogPath to not be set")
			}
		})
	}
}

func TestInstallMSIPackage(t *testing.T) {
	r, err := InstallMSIPackage("msi1", "Install", "{123}", "C:\\app.msi", "/quiet", "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != "Present" {
		t.Errorf("expected Ensure Present, got %v", r.Settings["Ensure"])
	}
}

func TestInstallMSIPackageEmptyPath(t *testing.T) {
	_, err := InstallMSIPackage("msi1", "Fail", "{123}", "", "", "", "", "")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestUninstallMSIPackage(t *testing.T) {
	r, err := UninstallMSIPackage("msi1", "Uninstall", "{123}", "C:\\app.msi", "", "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != "Absent" {
		t.Errorf("expected Ensure Absent, got %v", r.Settings["Ensure"])
	}
}

func TestUninstallMSIPackageEmptyPath(t *testing.T) {
	_, err := UninstallMSIPackage("msi1", "Fail", "{123}", "", "", "", "", "")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestMSIHashAlgorithmConstants(t *testing.T) {
	expected := map[string]string{
		"MD5":       FileHashMD5,
		"RIPEMD160": FileHashRIPEMD160,
		"SHA1":      FileHashSHA1,
		"SHA256":    FileHashSHA256,
		"SHA384":    FileHashSHA384,
		"SHA512":    FileHashSHA512,
	}
	for name, val := range expected {
		if val != name {
			t.Errorf("expected %q constant to be %q, got %q", name, name, val)
		}
	}
}
