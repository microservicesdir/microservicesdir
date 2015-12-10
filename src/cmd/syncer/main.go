package main

import (
	"flag"
	"os"
	"syncer"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
)

type project interface {
	Sync(id int, doneChannel chan int, errorsChannel chan error)
	IsAlreadyCheckedOut() (bool, error)
}

var (
	organization = flag.String("organization", "microservicesdir",
		"The Name of the organization you want to sync ex: github")
	typesOfRepos = flag.String("types", "all",
		"Type of Repositories you want to sync. [all, public, private]")
	target = flag.String("target", "target/repos", "Directory to store synced repos. [target/repos]")
)

func createTargetDirectory(targetDirectory string) error {
	_, err := os.Stat(targetDirectory)

	if os.IsNotExist(err) {
		log.Infof("Creating target directory: %v", *target)
		return os.MkdirAll(*target, os.FileMode(0777))
	}

	return err
}

func main() {
	flag.Parse()

	err := createTargetDirectory(*target)
	if err != nil {
		log.Fatalf("Could not create the target directory %v: %v", *target, err)
	}

	client := github.NewClient(nil)
	opt := &github.RepositoryListByOrgOptions{Type: *typesOfRepos}

	repos, _, err := client.Repositories.ListByOrg(*organization, opt)

	if err != nil {
		log.Errorf("Could not get the list of repositories for the organization %v. Error: %v", *organization, err)
	}

	if len(repos) == 0 {
		log.Infof("No repositories found for organization: %v", *organization)
		os.Exit(0)
	}

	numWorkers := len(repos)
	doneChan := make(chan int, numWorkers)
	errorsChannel := make(chan error)

	for i := 0; i <= numWorkers-1; i++ {
		var r = repos[i]

		var project = syncer.Project{
			Name:      *r.Name,
			GitURL:    *r.GitURL,
			TargetDir: *target,
		}

		go project.Sync(i, doneChan, errorsChannel)
	}

	doneCount := 0

	for doneCount < numWorkers {
		select {
		case e := <-errorsChannel:
			log.Error(e)
		case _ = <-doneChan:
			doneCount++
		}
	}

	log.Infof("Organization %v synced", *organization)
}
