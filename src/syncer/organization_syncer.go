package syncer

import (
	"core"

	log "github.com/Sirupsen/logrus"
)

// OrganizationSyncer knows how to sync the organization state with the local
// database
type OrganizationSyncer struct {
	Organization string
}

// SyncProjects sync all the organization projects with the local state
func (g *OrganizationSyncer) SyncProjects(rc RepositoriesClient) ([]core.Project, error) {
	// TODO: Github OAuth authentication to support private projects

	var projects []core.Project

	repos, err := rc.ListRepositories(g.Organization, "all")
	if err != nil {
		return nil, err
	}
	// TODO: persist the manifest information somewhere and update in case is newer.
	for _, v := range repos {
		manifest, err := rc.GetManifest("microservicesdir", "microservicesdir")
		if err != nil {
			log.Error(err)
		}

		log.Infof("%v", manifest)
		project := core.Project{
			Name: v.Name,
		}

		projects = append(projects, project)
	}

	return projects, err
}
