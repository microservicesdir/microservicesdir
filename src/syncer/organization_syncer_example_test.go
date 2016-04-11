package syncer_test

import (
	"net/http"
	"syncer"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/google/go-github/github"
)

func TestIntegrationSyncProjects(t *testing.T) {

	r, err := recorder.New("../../testdata/githubapi")

	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	client := &http.Client{
		Transport: r.Transport,
	}

	// TODO: Convert this into a example that can be attached to godocs
	o := syncer.OrganizationSyncer{
		Organization: "microservicesdir",
	}

	githubClient := github.NewClient(client)
	repositoriesClient := syncer.GithubRepositoryClient{githubClient}

	projects, err := o.SyncProjects(&repositoriesClient)

	if err != nil {
		t.Fatalf("error syncing projects %v", err)
	}

	if len(projects) < 1 {
		t.Error("Expected at least one project to be synced")
	}
}
