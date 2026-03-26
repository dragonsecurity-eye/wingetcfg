package wingetcfg

import "testing"

func TestNewLocalGroupResource(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		groupName     string
		description   string
		members       string
		ensure        string
		expectErr     bool
		expectEnsure  string
		expectMembers bool
	}{
		{
			name:          "Valid group with members",
			id:            "grp1",
			groupName:     "Administrators",
			description:   "Admin group",
			members:       "user1;user2",
			ensure:        EnsurePresent,
			expectErr:     false,
			expectEnsure:  EnsurePresent,
			expectMembers: true,
		},
		{
			name:          "Valid group without members",
			id:            "grp2",
			groupName:     "TestGroup",
			description:   "Test group",
			members:       "",
			ensure:        EnsurePresent,
			expectErr:     false,
			expectEnsure:  EnsurePresent,
			expectMembers: false,
		},
		{
			name:         "Remove group",
			id:           "grp3",
			groupName:    "TestGroup",
			description:  "",
			members:      "",
			ensure:       EnsureAbsent,
			expectErr:    false,
			expectEnsure: EnsureAbsent,
		},
		{
			name:      "Empty groupName returns error",
			id:        "grp4",
			groupName: "",
			ensure:    EnsurePresent,
			expectErr: true,
		},
		{
			name:         "Empty ID is allowed",
			id:           "",
			groupName:    "TestGroup",
			description:  "",
			members:      "",
			ensure:       EnsurePresent,
			expectErr:    false,
			expectEnsure: EnsurePresent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewLocalGroupResource(tc.id, tc.groupName, tc.description, tc.members, tc.ensure)
			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if r.Resource != WinGetLocalGroupResource {
				t.Errorf("expected resource %q, got %q", WinGetLocalGroupResource, r.Resource)
			}
			if tc.id != "" && r.ID != tc.id {
				t.Errorf("expected ID %q, got %q", tc.id, r.ID)
			}
			if tc.id == "" && r.ID != "" {
				t.Errorf("expected empty ID, got %q", r.ID)
			}
			if r.Settings["GroupName"] != tc.groupName {
				t.Errorf("expected GroupName %q, got %v", tc.groupName, r.Settings["GroupName"])
			}
			if r.Settings["Description"] != tc.description {
				t.Errorf("expected Description %q, got %v", tc.description, r.Settings["Description"])
			}
			if r.Settings["Ensure"] != tc.expectEnsure {
				t.Errorf("expected Ensure %q, got %v", tc.expectEnsure, r.Settings["Ensure"])
			}
			_, hasMembers := r.Settings["Members"]
			if tc.expectMembers && !hasMembers {
				t.Error("expected Members to be set")
			}
			if !tc.expectMembers && hasMembers {
				t.Error("expected Members to not be set")
			}
		})
	}
}

func TestAddOrModifyLocalGroup(t *testing.T) {
	r, err := AddOrModifyLocalGroup("grp1", "Admins", "Admin group", "user1;user2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsurePresent {
		t.Errorf("expected Ensure %q, got %v", EnsurePresent, r.Settings["Ensure"])
	}
	if r.Settings["GroupName"] != "Admins" {
		t.Errorf("expected GroupName %q, got %v", "Admins", r.Settings["GroupName"])
	}
	if r.Settings["Members"] != "user1;user2" {
		t.Errorf("expected Members %q, got %v", "user1;user2", r.Settings["Members"])
	}
}

func TestAddOrModifyLocalGroupEmptyGroupName(t *testing.T) {
	_, err := AddOrModifyLocalGroup("grp1", "", "Desc", "user1")
	if err == nil {
		t.Fatal("expected error for empty groupName, got nil")
	}
}

func TestIncludeMembersToGroup(t *testing.T) {
	r, err := IncludeMembersToGroup("grp1", "Admins", "user1;user2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsurePresent {
		t.Errorf("expected Ensure %q, got %v", EnsurePresent, r.Settings["Ensure"])
	}
	if r.Settings["GroupName"] != "Admins" {
		t.Errorf("expected GroupName %q, got %v", "Admins", r.Settings["GroupName"])
	}
	if r.Settings["MembersToInclude"] != "user1;user2" {
		t.Errorf("expected MembersToInclude %q, got %v", "user1;user2", r.Settings["MembersToInclude"])
	}
	if r.Resource != WinGetLocalGroupResource {
		t.Errorf("expected resource %q, got %q", WinGetLocalGroupResource, r.Resource)
	}
}

func TestIncludeMembersToGroupEmptyGroupName(t *testing.T) {
	_, err := IncludeMembersToGroup("grp1", "", "user1")
	if err == nil {
		t.Fatal("expected error for empty groupName, got nil")
	}
}

func TestIncludeMembersToGroupEmptyMembers(t *testing.T) {
	_, err := IncludeMembersToGroup("grp1", "Admins", "")
	if err == nil {
		t.Fatal("expected error for empty membersToInclude, got nil")
	}
}

func TestExcludeMembersFromGroup(t *testing.T) {
	r, err := ExcludeMembersFromGroup("grp1", "Admins", "user1;user2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsurePresent {
		t.Errorf("expected Ensure %q, got %v", EnsurePresent, r.Settings["Ensure"])
	}
	if r.Settings["GroupName"] != "Admins" {
		t.Errorf("expected GroupName %q, got %v", "Admins", r.Settings["GroupName"])
	}
	if r.Settings["MembersToExclude"] != "user1;user2" {
		t.Errorf("expected MembersToExclude %q, got %v", "user1;user2", r.Settings["MembersToExclude"])
	}
	if r.Resource != WinGetLocalGroupResource {
		t.Errorf("expected resource %q, got %q", WinGetLocalGroupResource, r.Resource)
	}
}

func TestExcludeMembersFromGroupEmptyGroupName(t *testing.T) {
	_, err := ExcludeMembersFromGroup("grp1", "", "user1")
	if err == nil {
		t.Fatal("expected error for empty groupName, got nil")
	}
}

func TestExcludeMembersFromGroupEmptyMembers(t *testing.T) {
	_, err := ExcludeMembersFromGroup("grp1", "Admins", "")
	if err == nil {
		t.Fatal("expected error for empty membersToExclude, got nil")
	}
}

func TestRemoveLocalGroup(t *testing.T) {
	r, err := RemoveLocalGroup("grp1", "TestGroup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Settings["Ensure"] != EnsureAbsent {
		t.Errorf("expected Ensure %q, got %v", EnsureAbsent, r.Settings["Ensure"])
	}
	if r.Settings["GroupName"] != "TestGroup" {
		t.Errorf("expected GroupName %q, got %v", "TestGroup", r.Settings["GroupName"])
	}
}

func TestRemoveLocalGroupEmptyGroupName(t *testing.T) {
	_, err := RemoveLocalGroup("grp1", "")
	if err == nil {
		t.Fatal("expected error for empty groupName, got nil")
	}
}

func TestIncludeMembersToGroupDirectives(t *testing.T) {
	r, err := IncludeMembersToGroup("grp1", "Admins", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Directives.Description != "Include members to group" {
		t.Errorf("expected description %q, got %q", "Include members to group", r.Directives.Description)
	}
	if !r.Directives.AllowPreRelease {
		t.Error("expected AllowPreRelease to be true")
	}
}

func TestExcludeMembersFromGroupDirectives(t *testing.T) {
	r, err := ExcludeMembersFromGroup("grp1", "Admins", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Directives.Description != "Exclude members from group" {
		t.Errorf("expected description %q, got %q", "Exclude members from group", r.Directives.Description)
	}
	if !r.Directives.AllowPreRelease {
		t.Error("expected AllowPreRelease to be true")
	}
}

func TestIncludeMembersToGroupEmptyID(t *testing.T) {
	r, err := IncludeMembersToGroup("", "Admins", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.ID != "" {
		t.Errorf("expected empty ID, got %q", r.ID)
	}
}

func TestExcludeMembersFromGroupEmptyID(t *testing.T) {
	r, err := ExcludeMembersFromGroup("", "Admins", "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.ID != "" {
		t.Errorf("expected empty ID, got %q", r.ID)
	}
}
