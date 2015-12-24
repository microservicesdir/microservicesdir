package core

import (
	"crypto/rand"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCanRetrieveAllProjects(t *testing.T) {
	t.Parallel()

	db, err := CreateDatabaseConnection("msvcdir", "msvcdir", "microservicesdirtest")
	if err != nil {
		t.Fatalf("could not connect to the database. %v", err)
	}

	defer db.Close()

	CreateProject(db)

	repo := ProjectRepository{db}
	projects := repo.GetAllProjects()

	if len(projects) < 1 {
		t.Errorf("there should be projects available. len %v want > 0", len(projects))
	}
}

// CreateProject will create a project using random information and return it.
// To be used in testing.
func CreateProject(db *sql.DB) Project {
	stmt, _ := db.Prepare("INSERT INTO projects (name, owner, language) VALUES(?,?,?)")
	defer stmt.Close()

	project := Project{
		Name:     randString(50),
		Language: randString(50),
		Owner:    randString(50),
	}

	stmt.Exec(project.Name, project.Owner, project.Language)
	return project
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
