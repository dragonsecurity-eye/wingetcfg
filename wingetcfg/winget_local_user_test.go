package wingetcfg

import "testing"

func TestNewLocalUserResource(t *testing.T) {
	tests := []struct {
		name                     string
		id                       string
		username                 string
		description              string
		disabled                 bool
		fullName                 string
		password                 string
		passwordChangeNotAllowed bool
		passwordChangeRequired   bool
		passwordNeverExpires     bool
		ensure                   string
		expectErr                bool
		expectEnsure             string
	}{
		{
			name:                     "Valid user with all fields",
			id:                       "user1",
			username:                 "testuser",
			description:              "Test user",
			disabled:                 false,
			fullName:                 "Test User",
			password:                 "P@ssw0rd",
			passwordChangeNotAllowed: true,
			passwordChangeRequired:   false,
			passwordNeverExpires:     true,
			ensure:                   EnsurePresent,
			expectErr:                false,
			expectEnsure:             EnsurePresent,
		},
		{
			name:         "Valid user minimal fields",
			id:           "",
			username:     "testuser",
			description:  "",
			disabled:     false,
			fullName:     "",
			password:     "",
			ensure:       EnsurePresent,
			expectErr:    false,
			expectEnsure: EnsurePresent,
		},
		{
			name:         "Remove user",
			id:           "user2",
			username:     "testuser",
			description:  "",
			disabled:     false,
			fullName:     "",
			password:     "",
			ensure:       EnsureAbsent,
			expectErr:    false,
			expectEnsure: EnsureAbsent,
		},
		{
			name:      "Empty username returns error",
			id:        "user3",
			username:  "",
			ensure:    EnsurePresent,
			expectErr: true,
		},
		{
			name:         "Disabled user",
			id:           "user4",
			username:     "disableduser",
			description:  "Disabled",
			disabled:     true,
			fullName:     "",
			password:     "",
			ensure:       EnsurePresent,
			expectErr:    false,
			expectEnsure: EnsurePresent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewLocalUserResource(tc.id, tc.username, tc.description, tc.disabled, tc.fullName, tc.password, tc.passwordChangeNotAllowed, tc.passwordChangeRequired, tc.passwordNeverExpires, tc.ensure)
			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if r.Resource != WinGetLocalUserResource {
				t.Errorf("expected resource %q, got %q", WinGetLocalUserResource, r.Resource)
			}
			if tc.id != "" && r.ID != tc.id {
				t.Errorf("expected ID %q, got %q", tc.id, r.ID)
			}
			if tc.id == "" && r.ID != "" {
				t.Errorf("expected empty ID, got %q", r.ID)
			}
			if r.Settings["UserName"] != tc.username {
				t.Errorf("expected UserName %q, got %v", tc.username, r.Settings["UserName"])
			}
			if r.Settings["Description"] != tc.description {
				t.Errorf("expected Description %q, got %v", tc.description, r.Settings["Description"])
			}
			if r.Settings["Disabled"] != tc.disabled {
				t.Errorf("expected Disabled %v, got %v", tc.disabled, r.Settings["Disabled"])
			}
			if r.Settings["FullName"] != tc.fullName {
				t.Errorf("expected FullName %q, got %v", tc.fullName, r.Settings["FullName"])
			}
			if tc.password != "" {
				if r.Settings["Password"] != tc.password {
					t.Errorf("expected Password %q, got %v", tc.password, r.Settings["Password"])
				}
			} else {
				if _, ok := r.Settings["Password"]; ok {
					t.Error("expected Password to not be set")
				}
			}
			if r.Settings["PasswordChangeNotAllowed"] != tc.passwordChangeNotAllowed {
				t.Errorf("expected PasswordChangeNotAllowed %v, got %v", tc.passwordChangeNotAllowed, r.Settings["PasswordChangeNotAllowed"])
			}
			if r.Settings["PasswordChangeRequired"] != tc.passwordChangeRequired {
				t.Errorf("expected PasswordChangeRequired %v, got %v", tc.passwordChangeRequired, r.Settings["PasswordChangeRequired"])
			}
			if r.Settings["PasswordNeverExpires"] != tc.passwordNeverExpires {
				t.Errorf("expected PasswordNeverExpires %v, got %v", tc.passwordNeverExpires, r.Settings["PasswordNeverExpires"])
			}
			if r.Settings["Ensure"] != tc.expectEnsure {
				t.Errorf("expected Ensure %q, got %v", tc.expectEnsure, r.Settings["Ensure"])
			}
		})
	}
}

func TestAddOrModifyLocalUser(t *testing.T) {
	r, err := AddOrModifyLocalUser("user1", "testuser", "Test user", false, "Test User", "P@ss", false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsurePresent {
		t.Errorf("expected Ensure %q, got %v", EnsurePresent, r.Settings["Ensure"])
	}
	if r.Settings["UserName"] != "testuser" {
		t.Errorf("expected UserName %q, got %v", "testuser", r.Settings["UserName"])
	}
}

func TestAddOrModifyLocalUserEmptyUsername(t *testing.T) {
	_, err := AddOrModifyLocalUser("user1", "", "Test", false, "", "", false, false, false)
	if err == nil {
		t.Fatal("expected error for empty username, got nil")
	}
}

func TestRemoveLocalUser(t *testing.T) {
	r, err := RemoveLocalUser("user1", "testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsureAbsent {
		t.Errorf("expected Ensure %q, got %v", EnsureAbsent, r.Settings["Ensure"])
	}
	if r.Settings["UserName"] != "testuser" {
		t.Errorf("expected UserName %q, got %v", "testuser", r.Settings["UserName"])
	}
}

func TestRemoveLocalUserEmptyUsername(t *testing.T) {
	_, err := RemoveLocalUser("user1", "")
	if err == nil {
		t.Fatal("expected error for empty username, got nil")
	}
}
