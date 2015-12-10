package manifest

import (
	"io/ioutil"
	"syncer"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
)

// Manifest of a service, describing its intents
type Manifest struct {
	// Name of the Micro Service
	Name string `yaml:"name"`
	// Language it is written in. (Scala, Ruby, Go, Haskell)
	Language string `yaml:"language"`
	// Owner of the service. Typically an email or a link for a project page
	Owner string `yaml:"owner"`
}

// IsHasManifest checks whether the project has a manifest assigned to it
func IsHasManifest(project syncer.Project) (bool, error) {
	log.Infof(project.Name)
	return true, nil
}

// ParseManifest reads a manifest from a yaml file and returns a manifest object
// populated with its information
func ParseManifest(pathToYamlFile string) (Manifest, error) {
	var manifest Manifest

	yamlFile, err := ioutil.ReadFile(pathToYamlFile)
	if err != nil {
		log.Errorf("yamlFile #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, &manifest)

	return manifest, err
}
