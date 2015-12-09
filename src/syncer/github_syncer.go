package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
)

//Project defines a git based project
type Project struct {
	name   string
	gitURL string
}

var (
	organization = flag.String("organization", "microservicesdir",
		"The name of the organization you want to sync ex: github")
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
	var wg sync.WaitGroup
	flag.Parse()

	err := createTargetDirectory(*target)
	if err != nil {
		log.Fatalf("Could not create the target directory %v: %v", *target, err)
	}

	client := github.NewClient(nil)
	opt := &github.RepositoryListByOrgOptions{Type: *typesOfRepos}

	repos, _, err := client.Repositories.ListByOrg(*organization, opt)

	if err != nil {
		log.Infof("Could not get the list of repositories for the organization %v", *organization)
		log.Error(err)
	}

	if len(repos) == 0 {
		log.Infof("No repositories found for organization: %v", *organization)
		os.Exit(0)
	}

	wg.Add(len(repos))
	errorsChannel := make(chan error)

	for _, r := range repos {
		var project = Project{
			name:   *r.Name,
			gitURL: *r.GitURL,
		}

		go project.sync(&wg, errorsChannel)
	}

	select {
	case err := <-errorsChannel:
		log.Error(err)
	default:
	}

	wg.Wait()
}

func (p Project) sync(wg *sync.WaitGroup, errors chan error) {
	isCheckedOut, err := p.isAlreadyCheckedOut()

	if err != nil {
		log.Error(err)
	}

	if isCheckedOut {
		errors <- p.update(wg)
	}

	errors <- p.clone(wg)
}

func (p Project) isAlreadyCheckedOut() (val bool, err error) {
	_, err = os.Stat(p.rootDir())
	switch {
	case err == nil:
		return true, nil
	case os.IsNotExist(err):
		return false, nil
	}

	return false, err
}

func (p Project) update(wg *sync.WaitGroup) error {
	defer wg.Done()

	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = p.rootDir()
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Couldn't update the repository %v: %v", p.name, err.Error())
	}

	cmd = exec.Command("git", "reset", "--hard", "origin/master")
	cmd.Dir = p.rootDir()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not reset the project %v", err)
	}

	log.Infof("Project %v updated", p.name)
	return nil
}

func (p Project) clone(wg *sync.WaitGroup) error {
	defer wg.Done()

	args := []string{"clone", p.gitURL, p.rootDir()}

	err := exec.Command("git", args...).Run()
	if err != nil {
		log.Errorf("Couldn't checkout the repository %v", p.name)
		log.Error(err)
		return err
	}

	log.Infof("Project %v was successfully cloned in %v", p.gitURL, p.rootDir())

	return nil
}

func (p Project) rootDir() string {
	return fmt.Sprintf("%v/%v", *target, p.name)
}
