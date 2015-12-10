package syncer

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
	Name   string
	GitURL string
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

// CreateProjectFromCheckout initializes a project based on an existing checkout
func CreateProjectFromCheckout(name string) Project {
	return Project{
		Name:   name,
		GitURL: name,
	}
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
		log.Errorf("Could not get the list of repositories for the organization %v. Error: %v", *organization, err)
	}

	if len(repos) == 0 {
		log.Infof("No repositories found for organization: %v", *organization)
		os.Exit(0)
	}

	wg.Add(len(repos))
	errorsChannel := make(chan error)

	for _, r := range repos {
		var project = Project{
			Name:   *r.Name,
			GitURL: *r.GitURL,
		}

		go project.sync(&wg, errorsChannel)
	}

	for err := range errorsChannel {
		log.Error(err)
	}

	wg.Wait()

}

func (p Project) sync(wg *sync.WaitGroup, errors chan error) {
	isCheckedOut, err := p.isAlreadyCheckedOut()

	if err != nil {
		log.Error(err)
	}

	if isCheckedOut {
		if err := p.update(wg); err != nil {
			errors <- err
		}
	}

	if err := p.clone(wg); err != nil {
		errors <- err
	}
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
		return fmt.Errorf("Couldn't update the repository %v: %v", p.Name, err.Error())
	}

	cmd = exec.Command("git", "reset", "--hard", "origin/master")
	cmd.Dir = p.rootDir()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not reset the project %v", err)
	}

	log.Infof("Project %v updated", p.Name)
	return nil
}

func (p Project) clone(wg *sync.WaitGroup) error {
	defer wg.Done()

	args := []string{"clone", p.GitURL, p.rootDir()}

	err := exec.Command("git", args...).Run()
	if err != nil {
		log.Errorf("Couldn't checkout the repository %v", p.Name)
		log.Error(err)
		return err
	}

	log.Infof("Project %v was successfully cloned in %v", p.GitURL, p.rootDir())

	return nil
}

func (p Project) rootDir() string {
	return fmt.Sprintf("%v/%v", *target, p.Name)
}
