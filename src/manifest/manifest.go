package manifest

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"syncer"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
)

// Manifest of a service, describing its intents
type Manifest struct {
	// Name of the Micro Service
	Name string `yaml:"name" json:"name"`
	// Language it is written in. (Scala, Ruby, Go, Haskell)
	Language string `yaml:"language" json:"language"`
	// Owner of the service. Typically an email or a link for a project page
	Owner string `yaml:"owner" json:"owner"`
}

// ToJSON returns a string json representation of the manifest
func (m *Manifest) ToJSON() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(b)
}

// FromJSON receives a JSON representation and returns the manifest associated
func FromJSON(data io.Reader) *Manifest {
	decoder := json.NewDecoder(data)
	var m Manifest

	err := decoder.Decode(&m)
	if err == nil {
		return &m
	}

	return nil
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
