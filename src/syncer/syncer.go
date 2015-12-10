package syncer

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

//Project defines a git based project
type Project struct {
	GitURL    string
	Name      string
	TargetDir string
}

// IsAlreadyCheckedOut returns true if the project is already checked and false otherwise
func (p Project) IsAlreadyCheckedOut() (val bool, err error) {
	_, err = os.Stat(p.rootDir())
	switch {
	case err == nil:
		return true, nil
	case os.IsNotExist(err):
		return false, nil
	}

	return false, err
}

// Sync syncs the project, cloning or updating it if necessary
func (p Project) Sync(id int, doneChannel chan int, errors chan error) {
	isCheckedOut, err := p.IsAlreadyCheckedOut()

	if err != nil {
		log.Error(err)
	}

	if isCheckedOut {
		if err := p.update(); err != nil {
			errors <- err
		}
	} else if err := p.clone(); err != nil {
		errors <- err
	}

	doneChannel <- id
}

func (p Project) update() error {
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = p.rootDir()
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Couldn't update the repository %v: %v", p.Name, err.Error())
	}

	cmd = exec.Command("git", "reset", "--hard", "origin")
	cmd.Dir = p.rootDir()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not reset the project %v. Error %v", p.Name, err)
	}

	log.Infof("Project %v updated", p.Name)
	return nil
}

func (p Project) clone() error {
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
	return fmt.Sprintf("%v/%v", p.TargetDir, p.Name)
}
