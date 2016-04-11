package syncer

import (
	"core"
	"net/http"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"

	"github.com/google/go-github/github"
)

// GithubRepositoryClient is an implementation of a Repository Client for github
type GithubRepositoryClient struct {
	*github.Client
}

// RepositoriesClient defines a client that knows how to interact with
type RepositoriesClient interface {
	ListRepositories(string, string) ([]core.Repository, error)
	GetManifest(string, string) (core.Manifest, error)
}

// GetManifest returns a string representation of the manifest file of a given repository.
// It will return a blank string in case there is no manifest for the project or an error.
func (gc *GithubRepositoryClient) GetManifest(organization string, repositoryName string) (core.Manifest, error) {
	var manifest core.Manifest

	fileContent, _, response, err := gc.Client.Repositories.GetContents(organization, repositoryName, "manifest.yml", nil)

	if response.StatusCode == http.StatusNotFound {
		return core.Manifest{}, err
	} else if err != nil {
		log.Fatalf("Error fetching manifest content. %v", err)
	}

	sDec, err := fileContent.Decode()
	if err != nil {
		log.Fatalf("Error decoding github content response. %v", err)
	}

	err = yaml.Unmarshal([]byte(sDec), &manifest)
	return manifest, err
}

// ListRepositories lists all the available repositories in an organization.
func (gc *GithubRepositoryClient) ListRepositories(organization string, typesOfRepos string) ([]core.Repository, error) {
	var repositories []core.Repository

	opt := &github.RepositoryListByOrgOptions{Type: typesOfRepos}
	repos, _, err := gc.Client.Repositories.ListByOrg(organization, opt)

	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		repository := core.Repository{
			Name:   *repo.Name,
			GitURL: *repo.GitURL,
		}
		repositories = append(repositories, repository)
	}

	return repositories, nil
}
