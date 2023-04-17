package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	// "github.com/elastic/go-elasticsearch/v7"
	// "github.com/elastic/go-elasticsearch/v7/esapi"
	// "github.com/elastic/go-elasticsearch/v7/estransport"
	// "github.com/elastic/go-elasticsearch/v7/esutil"
)

var Config config

func LoadConfig() error {
	var c config
	ymlFile, err := c.GetConfigFile()
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(ymlFile, &Config)
	if err != nil {
		return err
	}
	Config.SetDefault()
	return err
}

func (c *config) GetConfigFile() ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	settingsFile := "/internal/config/app_secret/doc_reco_settings.yml"
	fName := filepath.Join(dir, settingsFile)
	ymlFile, err := ioutil.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	return ymlFile, nil
}
