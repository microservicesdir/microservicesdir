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

	projects = append(projects, core.Project{Name: "microservicesdir", Language: "go", Owner: "vitorp@gmail.com"})
	projects = append(projects, core.Project{Name: "project2", Language: "go", Owner: "vitorp@gmail.com"})
	projects = append(projects, core.Project{Name: "project3", Language: "go", Owner: "vitorp@gmail.com"})

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
	expectedLanguage := mockDatabase.GetAllProjects()[0].Language
	expectedOwner := mockDatabase.GetAllProjects()[0].Owner

	returnedName := returnedProjects[0].Name
	returnedLanguage := returnedProjects[0].Language
	returnedOwner := returnedProjects[0].Owner

	if returnedName != expectedName {
		t.Errorf("project name should be %v but was %v", expectedName, returnedName)
	}

	if returnedLanguage != expectedLanguage {
		t.Errorf("project name should be %v but was %v", returnedLanguage, expectedLanguage)
	}
	if returnedOwner != expectedOwner {
		t.Errorf("project name should be %v but was %v", expectedOwner, returnedOwner)
	}
}
