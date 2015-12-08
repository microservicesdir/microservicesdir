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

//Project defines a git based project.
type Project struct {
	name   string
	gitURL string
}

var (
	wg sync.WaitGroup

	organization = flag.String("organization", "microservicesdir",
		"The name of the organization you want to sync ex: github")
	typesOfRepos = flag.String("types", "all",
		"Type of Repositories you want to sync. [all, public, private]")
	target = flag.String("target", "target/repos", "Directory to store synced repos. [target/repos]")
)

func main() {
	flag.Parse()

	log.Infof("Creating target folder: %v", *target)

	err := os.MkdirAll(*target, os.FileMode(0777))
	if err != nil {
		log.Fatalf("Could not create the target folder %v: %v", *target, err)
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

	for _, r := range repos {
		wg.Add(1)
		go syncRepository(Project{
			name:   *r.Name,
			gitURL: *r.GitURL,
		})
	}

	wg.Wait()
}

func syncRepository(project Project) {
	_, err := os.Stat(project.rootDir())
	if err == nil {
		updateProject(project)
	} else if os.IsNotExist(err) {
		cloneProject(project)
	}
}

func updateProject(project Project) {
	defer wg.Done()

	log.Infof("About to update project %v", project.name)
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = project.rootDir()
	err := cmd.Run()

	if err != nil {
		log.Infof("Couldn't update the repository %v", project.name)
		log.Error(err)
		return
	}

	cmd = exec.Command("git", "reset", "--hard", "origin/master")
	cmd.Dir = project.rootDir()
	_ = cmd.Run()

	log.Infof("Project %v updated", project.name)
}

func cloneProject(project Project) {
	defer wg.Done()

	log.Infof("About to clone %v in %v", project.gitURL, project.rootDir())
	args := []string{"clone", project.gitURL, project.rootDir()}

	err := exec.Command("git", args...).Run()
	if err != nil {
		log.Errorf("Couldn't checkout the repository %v", project.name)
		log.Error(err)
	}
}

func (p Project) rootDir() string {
	return fmt.Sprintf("%v/%v", target, p.name)
}
