package helm

import (
	"fmt"
	"net/url"

	helmrepo "helm.sh/helm/v3/pkg/repo"
)

// Repo holds resolved repository connection details.
type Repo struct {
	URL      string
	Username string
	Password string
}

// GetRepoByName looks up a configured Helm repository by name.
// repoFile is the path to repositories.yaml (e.g. from settings.RepositoryConfig).
func GetRepoByName(name, repoFile string) (*Repo, error) {
	f, err := helmrepo.LoadFile(repoFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load repository config: %w", err)
	}

	for _, entry := range f.Repositories {
		if entry.Name == name {
			return &Repo{
				URL:      entry.URL,
				Username: entry.Username,
				Password: entry.Password,
			}, nil
		}
	}
	return nil, fmt.Errorf("repo %q not found", name)
}

// RepoFromURL parses a raw URL, extracting any embedded credentials.
func RepoFromURL(rawURL string) (*Repo, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid repo URL %q: %w", rawURL, err)
	}

	repo := &Repo{}
	if u.User != nil {
		repo.Username = u.User.Username()
		repo.Password, _ = u.User.Password()
		u.User = nil
	}
	repo.URL = u.String()
	return repo, nil
}
