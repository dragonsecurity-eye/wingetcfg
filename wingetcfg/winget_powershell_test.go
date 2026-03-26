package wingetcfg

import "testing"

func TestExecutePowershellScript(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		scriptName string
		pwshell    string
		run        string
	}{
		{
			name:       "Basic script",
			id:         "ps1",
			scriptName: "TestScript",
			pwshell:    "Write-Host 'Hello'",
			run:        "once",
		},
		{
			name:       "Empty fields are allowed",
			id:         "",
			scriptName: "",
			pwshell:    "",
			run:        "",
		},
		{
			name:       "Complex script",
			id:         "ps2",
			scriptName: "InstallFeature",
			pwshell:    "Install-WindowsFeature -Name Web-Server -IncludeManagementTools",
			run:        "always",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := ExecutePowershellScript(tc.id, tc.scriptName, tc.pwshell, tc.run)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if r.Resource != OpenUEMPowershell {
				t.Errorf("expected resource %q, got %q", OpenUEMPowershell, r.Resource)
			}
			if r.Settings["ID"] != tc.id {
				t.Errorf("expected ID %q, got %v", tc.id, r.Settings["ID"])
			}
			if r.Settings["Name"] != tc.scriptName {
				t.Errorf("expected Name %q, got %v", tc.scriptName, r.Settings["Name"])
			}
			if r.Settings["Script"] != tc.pwshell {
				t.Errorf("expected Script %q, got %v", tc.pwshell, r.Settings["Script"])
			}
			if r.Settings["ScriptRun"] != tc.run {
				t.Errorf("expected ScriptRun %q, got %v", tc.run, r.Settings["ScriptRun"])
			}
		})
	}
}

func TestExecutePowershellScriptSettingsMapInitialized(t *testing.T) {
	r, err := ExecutePowershellScript("ps1", "Test", "echo hi", "once")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings == nil {
		t.Fatal("expected Settings map to be initialized, got nil")
	}
	if len(r.Settings) != 4 {
		t.Errorf("expected 4 settings, got %d", len(r.Settings))
	}
}

func TestExecutePowershellScriptNoError(t *testing.T) {
	// ExecutePowershellScript always returns nil error
	_, err := ExecutePowershellScript("", "", "", "")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestOpenUEMPowershellConstant(t *testing.T) {
	if OpenUEMPowershell != "openuem/Powershell" {
		t.Errorf("unexpected OpenUEMPowershell value: %q", OpenUEMPowershell)
	}
}
