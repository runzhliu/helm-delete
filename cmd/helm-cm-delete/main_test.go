package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetFieldsFromEnv(t *testing.T) {
	// Save the original environment values to restore later
	origUsername := os.Getenv("HELM_REPO_USERNAME")
	origPassword := os.Getenv("HELM_REPO_PASSWORD")
	origAccessToken := os.Getenv("HELM_REPO_ACCESS_TOKEN")
	origAuthHeader := os.Getenv("HELM_REPO_AUTH_HEADER")
	origContextPath := os.Getenv("HELM_REPO_CONTEXT_PATH")

	defer func() {
		_ = os.Setenv("HELM_REPO_USERNAME", origUsername)
		_ = os.Setenv("HELM_REPO_PASSWORD", origPassword)
		_ = os.Setenv("HELM_REPO_ACCESS_TOKEN", origAccessToken)
		_ = os.Setenv("HELM_REPO_AUTH_HEADER", origAuthHeader)
		_ = os.Setenv("HELM_REPO_CONTEXT_PATH", origContextPath)
	}()

	_ = os.Setenv("HELM_REPO_USERNAME", "testuser")
	_ = os.Setenv("HELM_REPO_PASSWORD", "testpass")
	_ = os.Setenv("HELM_REPO_ACCESS_TOKEN", "testtoken")
	_ = os.Setenv("HELM_REPO_AUTH_HEADER", "X-Test-Auth")
	_ = os.Setenv("HELM_REPO_CONTEXT_PATH", "/test/path")

	d := &deleteCmd{}
	d.setFieldsFromEnv()

	if d.username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", d.username)
	}
	if d.password != "testpass" {
		t.Errorf("Expected password 'testpass', got '%s'", d.password)
	}
	if d.accessToken != "testtoken" {
		t.Errorf("Expected accessToken 'testtoken', got '%s'", d.accessToken)
	}
	if d.authHeader != "X-Test-Auth" {
		t.Errorf("Expected authHeader 'X-Test-Auth', got '%s'", d.authHeader)
	}
	if d.contextPath != "/test/path" {
		t.Errorf("Expected contextPath '/test/path', got '%s'", d.contextPath)
	}
}

func TestResolveRepoURLWithDirectURL(t *testing.T) {
	d := &deleteCmd{}
	url, err := d.resolveRepoURL("https://admin:secret@chartmuseum.example.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if url != "https://chartmuseum.example.com" {
		t.Errorf("Expected URL 'https://chartmuseum.example.com', got '%s'", url)
	}
	if d.username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", d.username)
	}
	if d.password != "secret" {
		t.Errorf("Expected password 'secret', got '%s'", d.password)
	}
}

func TestResolveRepoURLWithExistingRepoConfig(t *testing.T) {
	// Create a temporary repositories.yaml file
	tempDir := t.TempDir()
	repoFile := filepath.Join(tempDir, "repositories.yaml")
	repoContent := `
apiVersion: v1
repositories:
- name: myrepo
  url: https://chartmuseum.mycompany.com
  username: repouser
  password: repopassword
`
	if err := os.WriteFile(repoFile, []byte(repoContent), 0644); err != nil {
		t.Fatalf("Failed to create temporary repo file: %v", err)
	}

	// Mock settings.RepositoryConfig by temporarily pointing HELM_REPOSITORY_CONFIG to our temp file.
	// We need to inject the mock, so we override the environment variable that helm/v3 uses under the hood.
	origHelmRepoConfig := os.Getenv("HELM_REPOSITORY_CONFIG")
	defer os.Setenv("HELM_REPOSITORY_CONFIG", origHelmRepoConfig)
	os.Setenv("HELM_REPOSITORY_CONFIG", repoFile)

	d := &deleteCmd{}
	url, err := d.resolveRepoURL("myrepo")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if url != "https://chartmuseum.mycompany.com" {
		t.Errorf("Expected URL 'https://chartmuseum.mycompany.com', got '%s'", url)
	}
	if d.username != "repouser" {
		t.Errorf("Expected username 'repouser', got '%s'", d.username)
	}
	if d.password != "repopassword" {
		t.Errorf("Expected password 'repopassword', got '%s'", d.password)
	}
}
