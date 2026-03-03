package helm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepoFromURL(t *testing.T) {
	tests := []struct {
		name        string
		rawURL      string
		wantURL     string
		wantUser    string
		wantPass    string
		expectError bool
	}{
		{
			name:        "URL without credentials",
			rawURL:      "https://chartmuseum.example.com",
			wantURL:     "https://chartmuseum.example.com",
			wantUser:    "",
			wantPass:    "",
			expectError: false,
		},
		{
			name:        "URL with credentials",
			rawURL:      "https://admin:secret123@chartmuseum.example.com/charts",
			wantURL:     "https://chartmuseum.example.com/charts",
			wantUser:    "admin",
			wantPass:    "secret123",
			expectError: false,
		},
		{
			name:        "URL with username only",
			rawURL:      "https://admin@chartmuseum.example.com",
			wantURL:     "https://chartmuseum.example.com",
			wantUser:    "admin",
			wantPass:    "",
			expectError: false,
		},
		{
			name:        "Invalid URL",
			rawURL:      "://bad-url",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := RepoFromURL(tt.rawURL)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError {
				if repo.URL != tt.wantURL {
					t.Errorf("expected URL %q, got %q", tt.wantURL, repo.URL)
				}
				if repo.Username != tt.wantUser {
					t.Errorf("expected Username %q, got %q", tt.wantUser, repo.Username)
				}
				if repo.Password != tt.wantPass {
					t.Errorf("expected Password %q, got %q", tt.wantPass, repo.Password)
				}
			}
		})
	}
}

func TestGetRepoByName(t *testing.T) {
	tempDir := t.TempDir()
	repoFile := filepath.Join(tempDir, "repositories.yaml")
	repoContent := `
apiVersion: v1
repositories:
- name: mytestrepo
  url: https://test.chartmuseum.com
  username: testuser
  password: testpassword
`
	if err := os.WriteFile(repoFile, []byte(repoContent), 0644); err != nil {
		t.Fatalf("Failed to write mock repo config: %v", err)
	}

	t.Run("Existing repo", func(t *testing.T) {
		repo, err := GetRepoByName("mytestrepo", repoFile)
		if err != nil {
			t.Fatalf("Unexpected error for existing repo: %v", err)
		}
		if repo.URL != "https://test.chartmuseum.com" {
			t.Errorf("Expected URL 'https://test.chartmuseum.com', got '%s'", repo.URL)
		}
		if repo.Username != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", repo.Username)
		}
		if repo.Password != "testpassword" {
			t.Errorf("Expected password 'testpassword', got '%s'", repo.Password)
		}
	})

	t.Run("Non-existent repo", func(t *testing.T) {
		_, err := GetRepoByName("doesnotexist", repoFile)
		if err == nil {
			t.Fatalf("Expected error for non-existent repo, got nil")
		}
	})
}
