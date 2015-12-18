package syncer

import log "github.com/Sirupsen/logrus"

// ProjectManager has the ability to fetch all the projects for a given organization.
// This interface abstract the details on where and how to do it.
type ProjectManager interface {
	GetAllProjects(*Organization) []Project
}

// AllProjects will query all available projects in the database
func AllProjects(lk Lookup) ([]*Project, error) {
	var projects []*Project
	db := lk.DB
	rows, err := db.Query("SELECT name, owner, language FROM projects")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		project := new(Project)
		err := rows.Scan(&project.Name, &project.Owner, &project.Language)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

// GithubProjectManager is a github implementation of a project manager
type GithubProjectManager struct{}

// Organization defines an organization having several different
// projects. Typically, this is the organzation where microservicesdir is
// deployed.
type Organization struct {
	Name      string
	GithubURL string

	ProjectManager
}

func (o *Organization) olha() {
	ps := o.ProjectManager.GetAllProjects(o)

	log.Info(ps[1])
}

// GetAllProjects gets all projects for a given organzation
func (ghm GithubProjectManager) GetAllProjects(o *Organization) []Project {
	var projects []Project

	projects[0] = Project{Name: "Hola que tal?"}

	return projects
}
