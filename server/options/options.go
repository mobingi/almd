package options

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Options struct {
	ID           string `yaml:"id"`
	Token        string `yaml:"token"`
	TemplateName string `yaml:"templateName"`
	ClusterName  string `yaml:"clusterName"`
}

// NewOptionsFromFile read config from disk,fileName must be a path so we can find it
func NewOptionsFromFile(fileName string) (*Options, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	o := &Options{}
	if err := yaml.Unmarshal(data, o); err != nil {
		return nil, err
	}

	return o, nil
}
