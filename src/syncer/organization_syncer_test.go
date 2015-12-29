package syncer

import (
	"core"
	"testing"
)

type fakeRepositoryManager struct{}

func (frm *fakeRepositoryManager) SaveProject(project core.Project, manifest core.Manifest) (core.Project, error) {
	return core.Project{}, nil
}

type fakeRepositoriesClient struct{}

func (frc *fakeRepositoriesClient) ListRepositories(organization string, types string) ([]core.Repository, error) {
	var repositories []core.Repository
	repository := core.Repository{
		Name:   "microservicesdir",
		GitURL: "git@github.com:microservicesdir/microservicesdir.git",
	}
	return append(repositories, repository), nil
}

func (frc *fakeRepositoriesClient) GetManifest(owner string, repositoryName string) (core.Manifest, error) {
	return core.Manifest{
		Name:     "microservicesdir",
		Owner:    "vitorp@gmail.com",
		Language: "go",
	}, nil
}

func TestCanParseManifest(t *testing.T) {
	fakeClient := fakeRepositoriesClient{}

	manifest, _ := fakeClient.GetManifest("microservicesdir", "microservicesdir")

	if manifest.Name != "microservicesdir" {
		t.Fatalf("cannot parse manifest. want %v got %v", "microservicesdir", manifest.Name)
	}
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
