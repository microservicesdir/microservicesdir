package core

// Project is a part of the organization's inventory of services.
type Project struct {
	Name string `json:"name"`
}

// Projects represents a list of projects
type Projects []Project

// GetAllProjects returns all the existing projects
func (p *Project) GetAllProjects() []Project {
	var projects []Project

	projects = append(projects, Project{Name: "microservicesdir"})
	projects = append(projects, Project{Name: "project2"})
	projects = append(projects, Project{Name: "project3"})

	return projects
}
