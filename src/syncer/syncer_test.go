package syncer_test

import (
	"database/sql"
	"fmt"
	"os"
	"syncer"
	"testing"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

const ()

func setupDB() (*sql.DB, error) {
	const (
		dbUser     = "msvcdir"
		dbPassword = "msvcdir"
		dbName     = "microservicesdirtest"
	)

	dbinfo := fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

func createTargetDirectory(targetDirectory string) error {
	_, err := os.Stat(targetDirectory)

	if os.IsNotExist(err) {
		return os.MkdirAll(targetDirectory, os.FileMode(0777))
	}

	return err
}

func TestSyncsTheOrganizationWithTheDatabase(t *testing.T) {
	t.Parallel()

	var projectName string

	project := setupProject()
	syncProject(project)

	db, err := setupDB()
	defer db.Close()

	if err != nil {
		t.Errorf("Could not connect to the database. %v", err)
		t.Fail()
	}

	row := db.QueryRow("select * from projects where name=? limit 1", project.Name)

	err = row.Scan(&projectName)

	switch {
	case err == sql.ErrNoRows:
		t.Logf("Expected a project named %v to be found in database, but it was not", project.Name)
		t.Fail()
	case err != nil:
		t.Log(err)
		t.Fail()
	}
}

func TestSyncsAnOrganizationInTheTargetDirectory(t *testing.T) {
	t.Parallel()

	project := setupProject()
	syncProject(project)

	projectCheckout := fmt.Sprintf("%v/%v", project.TargetDir, project.Name)
	_, err := os.Stat(projectCheckout)

	if os.IsNotExist(err) {
		t.Logf("Expected checkout directory %v for project %v to have been synced in the target directory", projectCheckout, project.Name)
		t.Fail()
	}
}

func syncProject(project syncer.Project) {
	numWorkers := 1
	doneChan := make(chan int, numWorkers)
	errorsChannel := make(chan error)
	doneCount := 0

	go project.Sync(1, doneChan, errorsChannel)

	for doneCount < numWorkers {
		select {
		case e := <-errorsChannel:
			log.Error(e)
		case _ = <-doneChan:
			doneCount++
		}
	}
}

func setupProject() syncer.Project {
	targetDirectory := "../../target/testing/repos"

	if err := createTargetDirectory(targetDirectory); err != nil {
		log.Errorf("Could not create target directory")
	}

	return syncer.Project{
		Name:      "microservicesdir",
		GitURL:    "http://github.com/microservicesdir/microservicesdir",
		TargetDir: targetDirectory,
	}
}
