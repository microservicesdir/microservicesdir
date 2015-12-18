package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"syncer"
)

// Project represents a project maintained by the organization
// It represents the ownership of a project.
type projectJSON struct {
	Name      string `json:"name"`
	GithubURL string `json:"githubUrl"`
	Owner     string `json:"owner"`
	Language  string `json:"language"`
	State     string `json:"state"`
}

func (p projectJSON) toJSON() []byte {
	s, _ := json.Marshal(p)
	return s
}

func dbConnection() *sql.DB {
	return nil
}

// TODO: Create handler to get all projects for an organization

// GetAllProjects returns a json with all projects
// This handler also expects a database connection to be provided to it by
// an injection pattern.
func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	ps, err := syncer.AllProjects(syncer.Lookup{
		DB: dbConnection(),
	})

	if err != nil {
		w.WriteHeader(500)
		return
	}

	p := projectJSON{
		Name:      ps[0].Name,
		GithubURL: "github.com/foo/bar",
		Owner:     ps[0].Owner,
		Language:  ps[0].Language,
		State:     "development",
	}

	w.WriteHeader(200)
	w.Write(p.toJSON())
}
