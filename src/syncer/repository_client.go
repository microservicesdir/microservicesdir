package syncer

import (
	"core"

	log "github.com/Sirupsen/logrus"

	"github.com/google/go-github/github"
)

// GithubRepositoryClient is an implementation of a Repository Client for github
type GithubRepositoryClient struct {
	github.Client
}

// RepositoriesClient defines a client that knows how to interact with
type RepositoriesClient interface {
	ListRepositories(string, string) ([]string, error)
	GetManifest(string, string) (core.Manifest, error)
}

// GetManifest returns a string representation of the manifest file of a given repository.
// It will return a blank string in case there is no manifest for the project or an error.
func (gc *GithubRepositoryClient) GetManifest(organization string, repositoryName string) (core.Manifest, error) {
	fileContent, directoryContent, response, err := gc.Client.Repositories.GetContents("microservicesdir", "microservicesdir", "README.md", nil)

	log.Infof("file %v directory %v response %v error %v", *fileContent.Content, directoryContent, response, err)
	return core.Manifest{}, nil
}

// ListRepositories is
func (gc *GithubRepositoryClient) ListRepositories(organization string, typesOfRepos string) ([]string, error) {
	var repositories []string

	opt := &github.RepositoryListByOrgOptions{Type: typesOfRepos}
	repos, _, err := gc.Client.Repositories.ListByOrg(organization, opt)

	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		repositories = append(repositories, *repo.Name)
	}

	return repositories, nil
}
