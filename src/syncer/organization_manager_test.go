package syncer

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"
)

type fakeProjectManager struct{}

func (f *fakeProjectManager) GetAllProjects(o *Organization) []Project {
	var project []Project

	return project
}

func OpenDbConnection() (*sql.DB, error) {
	dbinfo := "msvcdir:msvcdir@/microservicesdirtest"
	db, err := sql.Open("mysql", dbinfo)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestOrganizationCanReturnAllProjects(t *testing.T) {
	t.Parallel()

	dbConn, err := OpenDbConnection()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	lookup := Lookup{
		DB: dbConn,
	}

	stmt, err := dbConn.Prepare("INSERT INTO projects(name, owner, language) VALUES (?,?,?)")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer stmt.Close()

	s, err := stmt.Exec(randomString(5), "vitorp@gmail.com", "go")

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	rowsAffected, err := s.LastInsertId()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if rowsAffected < 1 {
		t.Error("Should have inserted a new project but it failed")
		t.Fail()
	}

	projects, err := AllProjects(lookup)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if len(projects) == 0 {
		t.Log("Should return more than one project")
		t.Fail()
	}
}
