package syncer_test

import (
	"syncer"
	"testing"

	"github.com/google/go-github/github"
)

func _TestIntegrationSyncProjects(t *testing.T) {
	o := syncer.OrganizationSyncer{
		Organization: "microservicesdir",
	}

	githubClient := github.NewClient(nil)
	repositoriesClient := syncer.GithubRepositoryClient{*githubClient}

	projects, err := o.SyncProjects(&repositoriesClient)

	if err != nil {
		t.Fatalf("error syncing projects %v", err)
	}

	if len(projects) < 1 {
		t.Error("Expected at least one project to be synced")
	}
}
