package local

import (
	"io/ioutil"

	"launchpad.net/goyaml"
)

// BoshConfig describes a local ~/.bosh_config file
// See testhelpers/fixtures/bosh_config.yml
type BoshConfig struct {
	Target         string
	Name           string `yaml:"target_name"`
	Version        string `yaml:"target_version"`
	UUID           string `yaml:"target_uuid"`
	Aliases        map[string]map[string]string
	Authentication map[string]authentication `yaml:"auth"`
}

type authentication struct {
	Username string
	Password string
}

// LoadBoshConfig loads and unmarshals ~/.bosh_config
func LoadBoshConfig(configPath string) (config *BoshConfig, err error) {
	config = &BoshConfig{}

	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	goyaml.Unmarshal(contents, config)
	return
}

// CurrentBoshTarget returns the connection information for local user's current target BOSH
func (config *BoshConfig) CurrentBoshTarget() (target, username, password string) {
	return "", "", ""
}
