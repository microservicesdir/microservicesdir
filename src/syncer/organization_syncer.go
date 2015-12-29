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

// SyncProjects sync all the organization projects with
func (g *OrganizationSyncer) SyncProjects(rc RepositoriesClient) ([]core.Project, error) {
	var projects []core.Project

	repos, err := rc.ListRepositories(g.Organization, "all")
	if err != nil {
		return nil, err
	}

	for _, v := range repos {
		manifest, err := rc.GetManifest("microservicesdir", "microservicesdir")
		if err != nil {
			log.Error(err)
		}

		log.Infof("%v", manifest)

		projects = append(projects, core.Project{
			Name: v,
		})
	}

	// Sync their state in the database.
	//https://github.com/microservicesdir/microservicesdir/blob/master/Makefile
	//https://raw.githubusercontent.com/microservicesdir/microservicesdir/master/Makefile
	return projects, err
}
