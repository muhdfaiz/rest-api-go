package systems

import (
	"io/ioutil"
	"os"

	"github.com/smallfish/simpleyaml"
)

// Configs struct
type Configs struct {
}

// Get value from config file using Simple YAML package
func (conf *Configs) Get(fileName string, key string, defaultValue string) string {
	helpers := &Helpers{}
	file := helpers.StrConcat(os.Getenv("GOPATH"), "src/bitbucket.org/shoppermate/configs/", fileName)

	// Read YAML Config File
	source, err := ioutil.ReadFile(file)

	if err != nil {
		panic(err)
	}

	yaml, err := simpleyaml.NewYaml(source)

	if err != nil {
		panic(err)
	}

	value, err := yaml.Get(key).String()
	// Return default value if cannot retrieve config
	if err != nil {
		return defaultValue
	}

	return value
}
