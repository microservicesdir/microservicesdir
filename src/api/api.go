package api

import (
	"core"
	"encoding/json"
	"net/http"
)

// MicroServicesDirectoryDatabase is a DAO containing all operations
// supported by this app.
type MicroServicesDirectoryDatabase interface {
	GetAllProjects() []core.Project
}

// HomeHandler returns a http handler that returns a json representation for
// the homepage
func HomeHandler(db MicroServicesDirectoryDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projects, _ := json.Marshal(db.GetAllProjects())

		w.Write(projects)
	})
}
