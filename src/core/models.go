package core

import (
	"database/sql"
	"log"
)

// Project is a part of the organization's inventory of services.
type Project struct {
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Language string `json:"language"`
}

// Projects represents a list of projects
type Projects []Project

// ProjectRepository manages all the projects
type ProjectRepository struct {
	*sql.DB
}

// GetAllProjects returns all the existing projects
func (p *ProjectRepository) GetAllProjects() []Project {
	var projects []Project

	rows, err := p.DB.Query("select name, language, owner from projects")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.Name, &project.Language, &project.Owner); err != nil {
			log.Fatal(err)
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return projects
}
