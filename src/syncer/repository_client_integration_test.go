package syncer_test

import (
	"core"
	"syncer"
	"testing"

	"github.com/google/go-github/github"
)

func TestManifestParsing(t *testing.T) {
	githubClient := github.NewClient(nil)
	repositoriesClient := syncer.GithubRepositoryClient{*githubClient}

	manifest, err := repositoriesClient.GetManifest("microservicesdir", "microservicesdir")

	if err != nil {
		t.Fatalf("error parsing manifest projects %v", err)
	}

	if manifest.Name != "microservicesdir" {
		t.Fatalf("manifest name: got %v want %v", manifest.Name, "microservicesdir")
	}

	if manifest.Owner != "vitorp@gmail.com" {
		t.Fatalf("manifest owner: got %v want %v", manifest.Owner, "vitorp@gmail.com")
	}

	if manifest.Language != "go" {
		t.Fatalf("manifest language: got %v want %v", manifest.Language, "go")
	}
}

func TestNoAvailableManifestShouldBeBlank(t *testing.T) {
	githubClient := github.NewClient(nil)
	repositoriesClient := syncer.GithubRepositoryClient{*githubClient}

	manifest, _ := repositoriesClient.GetManifest("microservicesdir", "inexistingproject")
	var blankManifest core.Manifest

	if manifest != blankManifest {
		t.Fatalf("want a blank manifest but got %v", manifest)
	}
}
