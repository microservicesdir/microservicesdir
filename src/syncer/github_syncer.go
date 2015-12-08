package main

import (
	"flag"
	"os"
	"os/exec"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
)

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

	log.Info("Creating target folder: " + *target)

	err := os.MkdirAll(*target, os.FileMode(0777))
	if err != nil {
		log.Info("Could not create the target folder " + *target)
		log.Error(err)
		panic("Exiting")
	}

	client := github.NewClient(nil)
	opt := &github.RepositoryListByOrgOptions{Type: *typesOfRepos}

	repos, _, err := client.Repositories.ListByOrg(*organization, opt)

	if err != nil {
		log.Info("Could not get the list of repositories for the organization " + *organization)
		log.Error(err)
	} else if len(repos) == 0 {
		log.Info("No repositories found for organization: " + *organization)
		os.Exit(0)
	}

	for _, r := range repos {
		wg.Add(1)
		go syncRepository(*r.Name, *r.GitURL)
	}

	wg.Wait()
}

func syncRepository(name string, gitURL string) {
	_, err := os.Stat(projectRootDir(name))
	if err == nil {
		updateProject(name)
	} else if os.IsNotExist(err) {
		cloneProject(name, gitURL)
	}
}

func updateProject(name string) {
	defer wg.Done()

	log.Info("About to update project " + name)
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = projectRootDir(name)
	err := cmd.Run()

	if err != nil {
		log.Info("Couldn't update the repository " + name)
		log.Error(err)
	} else {
		cmd := exec.Command("git", "reset", "--hard", "origin/master")
		cmd.Dir = projectRootDir(name)
		_ = cmd.Run()

		log.Info("Project " + name + " updated")
	}
}

func cloneProject(name string, gitURL string) {
	defer wg.Done()

	log.Info("About to clone " + gitURL + " in  " + projectRootDir(name))
	args := []string{"clone", gitURL, projectRootDir(name)}

	err := exec.Command("git", args...).Run()
	if err != nil {
		log.Info("Couldn't checkout the repository " + name)
		log.Error(err)
	}
}

func projectRootDir(projectName string) string {
	return *target + "/" + projectName
}
