package wingetcfg

import "testing"

func TestIsValidRegistryValueType(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"String", true},
		{"Binary", true},
		{"DWord", true},
		{"QWord", true},
		{"MultiString", true},
		{"ExpandString", true},
		{"string", false},
		{"dword", false},
		{"Invalid", false},
		{"", false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := IsValidRegistryValueType(tc.input)
			if result != tc.expected {
				t.Errorf("IsValidRegistryValueType(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestNewWinGetRegistryResource(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		description string
		key         string
		valueName   string
		valueType   string
		valueData   string
		ensure      string
		hex         bool
		force       bool
		expectErr   bool
		errMsg      string
	}{
		{
			name:        "Valid key only",
			id:          "reg1",
			description: "Add key",
			key:         "HKLM:\\SOFTWARE\\Test",
			valueName:   "",
			valueType:   "",
			valueData:   "",
			ensure:      EnsurePresent,
			hex:         false,
			force:       false,
			expectErr:   false,
		},
		{
			name:        "Valid with value",
			id:          "reg2",
			description: "Add value",
			key:         "HKLM:\\SOFTWARE\\Test",
			valueName:   "TestValue",
			valueType:   "String",
			valueData:   "hello",
			ensure:      EnsurePresent,
			hex:         false,
			force:       false,
			expectErr:   false,
		},
		{
			name:        "DWord with hex",
			id:          "reg3",
			description: "DWord hex",
			key:         "HKLM:\\SOFTWARE\\Test",
			valueName:   "NumVal",
			valueType:   "DWord",
			valueData:   "0xFF",
			ensure:      EnsurePresent,
			hex:         true,
			force:       false,
			expectErr:   false,
		},
		{
			name:        "QWord with hex",
			id:          "reg4",
			description: "QWord hex",
			key:         "HKLM:\\SOFTWARE\\Test",
			valueName:   "BigVal",
			valueType:   "QWord",
			valueData:   "0xFFFF",
			ensure:      EnsurePresent,
			hex:         true,
			force:       false,
			expectErr:   false,
		},
		{
			name:        "Empty key returns error",
			id:          "reg5",
			description: "Fail",
			key:         "",
			valueName:   "",
			valueType:   "",
			valueData:   "",
			ensure:      EnsurePresent,
			hex:         false,
			force:       false,
			expectErr:   true,
			errMsg:      "key cannot be empty",
		},
		{
			name:        "Invalid value type returns error",
			id:          "reg6",
			description: "Fail",
			key:         "HKLM:\\SOFTWARE\\Test",
			valueName:   "TestValue",
			valueType:   "InvalidType",
			valueData:   "hello",
			ensure:      EnsurePresent,
			hex:         false,
			force:       false,
			expectErr:   true,
			errMsg:      "value type is not valid",
		},
		{
			name:        "Multiline valueData returns error",
			id:          "reg7",
			description: "Fail",
			key:         "HKLM:\\SOFTWARE\\Test",
			valueName:   "TestValue",
			valueType:   "String",
			valueData:   "line1\nline2",
			ensure:      EnsurePresent,
			hex:         false,
			force:       false,
			expectErr:   true,
		},
		{
			name:        "Force flag is set",
			id:          "reg8",
			description: "Force",
			key:         "HKLM:\\SOFTWARE\\Test",
			valueName:   "",
			valueType:   "",
			valueData:   "",
			ensure:      EnsurePresent,
			hex:         false,
			force:       true,
			expectErr:   false,
		},
		{
			name:        "Empty ID is allowed",
			id:          "",
			description: "No ID",
			key:         "HKLM:\\SOFTWARE\\Test",
			valueName:   "",
			valueType:   "",
			valueData:   "",
			ensure:      EnsurePresent,
			hex:         false,
			force:       false,
			expectErr:   false,
		},
		{
			name:        "Invalid ensure defaults to Present",
			id:          "reg9",
			description: "Bad ensure",
			key:         "HKLM:\\SOFTWARE\\Test",
			valueName:   "",
			valueType:   "",
			valueData:   "",
			ensure:      "Invalid",
			hex:         false,
			force:       false,
			expectErr:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewWinGetRegistryResource(tc.id, tc.description, tc.key, tc.valueName, tc.valueType, tc.valueData, tc.ensure, tc.hex, tc.force)
			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tc.errMsg != "" && err.Error() != tc.errMsg {
					t.Errorf("expected error %q, got %q", tc.errMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if r.Resource != WinGetRegistryResource {
				t.Errorf("expected resource %q, got %q", WinGetRegistryResource, r.Resource)
			}
			if tc.id != "" && r.ID != tc.id {
				t.Errorf("expected ID %q, got %q", tc.id, r.ID)
			}
			if tc.id == "" && r.ID != "" {
				t.Errorf("expected empty ID, got %q", r.ID)
			}
			if r.Settings["Key"] != tc.key {
				t.Errorf("expected Key %q, got %v", tc.key, r.Settings["Key"])
			}
			if r.Settings["ValueName"] != tc.valueName {
				t.Errorf("expected ValueName %q, got %v", tc.valueName, r.Settings["ValueName"])
			}
			if tc.valueType != "" {
				if r.Settings["ValueType"] != tc.valueType {
					t.Errorf("expected ValueType %q, got %v", tc.valueType, r.Settings["ValueType"])
				}
			}
			if tc.force {
				if r.Settings["Force"] != true {
					t.Error("expected Force to be true")
				}
			} else {
				if _, ok := r.Settings["Force"]; ok {
					t.Error("expected Force to not be set")
				}
			}
			if tc.valueType == RegistryValueTypeDWord || tc.valueType == RegistryValueTypeQWord {
				if r.Settings["Hex"] != tc.hex {
					t.Errorf("expected Hex %v, got %v", tc.hex, r.Settings["Hex"])
				}
			}
		})
	}
}

func TestAddRegistryKey(t *testing.T) {
	r, err := AddRegistryKey("reg1", "Add key", "HKLM:\\SOFTWARE\\Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsurePresent {
		t.Errorf("expected Ensure %q, got %v", EnsurePresent, r.Settings["Ensure"])
	}
	if r.Settings["Key"] != "HKLM:\\SOFTWARE\\Test" {
		t.Errorf("expected Key %q, got %v", "HKLM:\\SOFTWARE\\Test", r.Settings["Key"])
	}
}

func TestAddRegistryKeyEmptyKey(t *testing.T) {
	_, err := AddRegistryKey("reg1", "Fail", "")
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestAddRegistryValue(t *testing.T) {
	r, err := AddRegistryValue("reg1", "Add value", "HKLM:\\SOFTWARE\\Test", "MyValue", "String", "hello", false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsurePresent {
		t.Errorf("expected Ensure %q, got %v", EnsurePresent, r.Settings["Ensure"])
	}
	if r.Settings["ValueName"] != "MyValue" {
		t.Errorf("expected ValueName %q, got %v", "MyValue", r.Settings["ValueName"])
	}
	if r.Settings["ValueData"] != "hello" {
		t.Errorf("expected ValueData %q, got %v", "hello", r.Settings["ValueData"])
	}
}

func TestRemoveRegistryKey(t *testing.T) {
	r, err := RemoveRegistryKey("reg1", "Remove key", "HKLM:\\SOFTWARE\\Test", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsureAbsent {
		t.Errorf("expected Ensure %q, got %v", EnsureAbsent, r.Settings["Ensure"])
	}
	if r.Settings["Force"] != true {
		t.Error("expected Force to be true")
	}
}

func TestRemoveRegistryKeyEmptyKey(t *testing.T) {
	_, err := RemoveRegistryKey("reg1", "Fail", "", false)
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestRemoveRegistryValue(t *testing.T) {
	r, err := RemoveRegistryValue("reg1", "Remove value", "HKLM:\\SOFTWARE\\Test", "MyValue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsureAbsent {
		t.Errorf("expected Ensure %q, got %v", EnsureAbsent, r.Settings["Ensure"])
	}
	if r.Settings["ValueName"] != "MyValue" {
		t.Errorf("expected ValueName %q, got %v", "MyValue", r.Settings["ValueName"])
	}
}

func TestRemoveRegistryValueEmptyKey(t *testing.T) {
	_, err := RemoveRegistryValue("reg1", "Fail", "", "MyValue")
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestUpdateRegistryKeyDefaultValue(t *testing.T) {
	r, err := UpdateRegistryKeyDefaultValue("reg1", "Update default", "HKLM:\\SOFTWARE\\Test", "String", "default_val", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsurePresent {
		t.Errorf("expected Ensure %q, got %v", EnsurePresent, r.Settings["Ensure"])
	}
	if r.Settings["ValueData"] != "default_val" {
		t.Errorf("expected ValueData %q, got %v", "default_val", r.Settings["ValueData"])
	}
	if r.Settings["Force"] != true {
		t.Error("expected Force to be true")
	}
	// ValueName should be empty string for default value
	if r.Settings["ValueName"] != "" {
		t.Errorf("expected empty ValueName, got %v", r.Settings["ValueName"])
	}
}

func TestUpdateRegistryKeyDefaultValueInvalidType(t *testing.T) {
	_, err := UpdateRegistryKeyDefaultValue("reg1", "Fail", "HKLM:\\SOFTWARE\\Test", "BadType", "val", false)
	if err == nil {
		t.Fatal("expected error for invalid value type, got nil")
	}
}

func TestRegistryValueTypeConstants(t *testing.T) {
	if RegistryValueTypeString != "String" {
		t.Errorf("unexpected RegistryValueTypeString: %q", RegistryValueTypeString)
	}
	if RegistryValueTypeBinary != "Binary" {
		t.Errorf("unexpected RegistryValueTypeBinary: %q", RegistryValueTypeBinary)
	}
	if RegistryValueTypeDWord != "DWord" {
		t.Errorf("unexpected RegistryValueTypeDWord: %q", RegistryValueTypeDWord)
	}
	if RegistryValueTypeQWord != "QWord" {
		t.Errorf("unexpected RegistryValueTypeQWord: %q", RegistryValueTypeQWord)
	}
	if RegistryValueTypeMultistring != "MultiString" {
		t.Errorf("unexpected RegistryValueTypeMultistring: %q", RegistryValueTypeMultistring)
	}
	if RegistryValueTypeExpandString != "ExpandString" {
		t.Errorf("unexpected RegistryValueTypeExpandString: %q", RegistryValueTypeExpandString)
	}
}
