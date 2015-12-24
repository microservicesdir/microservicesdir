package api

import (
	"core"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockDb struct{}

func (db MockDb) GetAllProjects() []core.Project {
	return returnedProjects
}

var (
	returnedProjects = returnProjects()
	mockDatabase     MockDb
)

func returnProjects() []core.Project {
	var projects []core.Project

	projects = append(projects, core.Project{Name: "microservicesdir"})
	projects = append(projects, core.Project{Name: "project2"})
	projects = append(projects, core.Project{Name: "project3"})

	return projects
}

func TestReturnsAllProjectsAsJson(t *testing.T) {
	homeHandle := HomeHandler(mockDatabase)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)

	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}

	json.Unmarshal(body, &returnedProjects)
	expectedName := mockDatabase.GetAllProjects()[0].Name
	returnedName := returnedProjects[0].Name

	if returnedName != expectedName {
		t.Errorf("project name should be %v but was %v", expectedName, returnedName)
	}
}
