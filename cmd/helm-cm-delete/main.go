package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/runzhliu/helm-delete/pkg/chartmuseum"
	helmutil "github.com/runzhliu/helm-delete/pkg/helm"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/cli"
)

const defaultTimeout = 30

type deleteCmd struct {
	username    string
	password    string
	accessToken string
	authHeader  string
	contextPath string
	caFile      string
	certFile    string
	keyFile     string
	insecure    bool
	timeout     int64
}

func main() {
	cmd := newDeleteCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newDeleteCmd() *cobra.Command {
	d := &deleteCmd{}

	cmd := &cobra.Command{
		Use:   "helm cm-delete [NAME] [VERSION] [REPO]",
		Short: "Delete a chart version from ChartMuseum",
		Long: `Delete a specific version of a Helm chart from a ChartMuseum repository.

REPO may be a repository name (as configured via 'helm repo add') or a direct URL.

Examples:
  helm cm-delete mychart 1.2.3 myrepo
  helm cm-delete mychart 1.2.3 https://chartmuseum.example.com`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			d.setFieldsFromEnv()
			return d.delete(cmd, args[0], args[1], args[2])
		},
		SilenceUsage: true,
	}

	f := cmd.Flags()
	f.StringVarP(&d.username, "username", "u", "", "Override chart repository username")
	f.StringVarP(&d.password, "password", "p", "", "Override chart repository password")
	f.StringVar(&d.accessToken, "access-token", "", "Send token in Authorization header")
	f.StringVar(&d.authHeader, "auth-header", "", "Send token in custom header (e.g. X-Auth-Token)")
	f.StringVar(&d.contextPath, "context-path", "", "ChartMuseum context path (when behind a reverse proxy)")
	f.StringVar(&d.caFile, "ca-file", "", "Verify certificates using this CA bundle")
	f.StringVar(&d.certFile, "cert-file", "", "Identify HTTPS client using this SSL certificate file")
	f.StringVar(&d.keyFile, "key-file", "", "Identify HTTPS client using this SSL key file")
	f.BoolVarP(&d.insecure, "insecure", "i", false, "Skip TLS certificate verification")
	f.Int64Var(&d.timeout, "timeout", defaultTimeout, "Request timeout in seconds")

	return cmd
}

// setFieldsFromEnv populates unset fields from environment variables,
// mirroring the same env vars used by helm-push for consistency.
func (d *deleteCmd) setFieldsFromEnv() {
	if v, ok := os.LookupEnv("HELM_REPO_USERNAME"); ok && d.username == "" {
		d.username = v
	}
	if v, ok := os.LookupEnv("HELM_REPO_PASSWORD"); ok && d.password == "" {
		d.password = v
	}
	if v, ok := os.LookupEnv("HELM_REPO_ACCESS_TOKEN"); ok && d.accessToken == "" {
		d.accessToken = v
	}
	if v, ok := os.LookupEnv("HELM_REPO_AUTH_HEADER"); ok && d.authHeader == "" {
		d.authHeader = v
	}
	if v, ok := os.LookupEnv("HELM_REPO_CONTEXT_PATH"); ok && d.contextPath == "" {
		d.contextPath = v
	}
}

func (d *deleteCmd) delete(cmd *cobra.Command, name, version, repoArg string) error {
	repoURL, err := d.resolveRepoURL(repoArg)
	if err != nil {
		return err
	}

	client, err := chartmuseum.NewClient(
		chartmuseum.URL(repoURL),
		chartmuseum.Username(d.username),
		chartmuseum.Password(d.password),
		chartmuseum.AccessToken(d.accessToken),
		chartmuseum.AuthHeader(d.authHeader),
		chartmuseum.ContextPath(d.contextPath),
		chartmuseum.CAFile(d.caFile),
		chartmuseum.CertFile(d.certFile),
		chartmuseum.KeyFile(d.keyFile),
		chartmuseum.InsecureSkipVerify(d.insecure),
		chartmuseum.Timeout(d.timeout),
	)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Deleting %s-%s from %s...\n", name, version, repoURL)

	if err := client.DeleteChartVersion(name, version); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Successfully deleted %s-%s\n", name, version)
	return nil
}

// resolveRepoURL returns a usable URL for the given repo argument, which may be
// a configured repo name or a direct URL (http:// / https://).
func (d *deleteCmd) resolveRepoURL(repoArg string) (string, error) {
	if isURL(repoArg) {
		repo, err := helmutil.RepoFromURL(repoArg)
		if err != nil {
			return "", err
		}
		// Credentials embedded in the URL take the lowest priority.
		if d.username == "" {
			d.username = repo.Username
		}
		if d.password == "" {
			d.password = repo.Password
		}
		return repo.URL, nil
	}

	// Treat as a named Helm repository.
	settings := cli.New()
	repo, err := helmutil.GetRepoByName(repoArg, settings.RepositoryConfig)
	if err != nil {
		return "", err
	}
	// Repo-stored credentials are used as fallback when not overridden by flags/env.
	if d.username == "" {
		d.username = repo.Username
	}
	if d.password == "" {
		d.password = repo.Password
	}
	return repo.URL, nil
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
