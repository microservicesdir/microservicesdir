package syncer

import (
	"core"
	"testing"

	"gopkg.in/yaml.v2"
)

type fakeRepositoriesClient struct{}

func (frc *fakeRepositoriesClient) ListRepositories(organization string, types string) ([]string, error) {
	return []string{"microservicesdir", types}, nil
}

func (frc *fakeRepositoriesClient) GetManifest(owner string, repositoryName string) (core.Manifest, error) {
	manifestStr := `
name: ExampleProject
owner: developer-team@example.com
language: go
`
	var manifest core.Manifest
	err := yaml.Unmarshal([]byte(manifestStr), &manifest)

	return manifest, err
}

func TestCanRetrieveAllProjectsFromGithub(t *testing.T) {
	o := OrganizationSyncer{
		Organization: "microservicesdir",
	}

	fakeClient := fakeRepositoriesClient{}

	projects, err := o.SyncProjects(&fakeClient)

	if err != nil {
		t.Fatalf("error syncing projects %v", err)
	}

	if len(projects) < 1 {
		t.Error("Expected at least one project to be synced")
	}
}
