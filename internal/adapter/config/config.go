package config

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/port"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"gopkg.in/yaml.v3"
	"os"
)

const appDir = "/etc/aictl"
const configPath = appDir + "/context.yaml"

var _ port.Config = &Adapter{}

type Adapter struct {
}

func NewContextAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) GetContextFromAictlFolder() *config.Config {
	if !fshelper.PathExists(appDir) || !fshelper.PathExists(configPath) {
		return &config.Config{}
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		// TODO add log

		return &config.Config{}
	}

	var fileCfg fileConfig
	err = yaml.Unmarshal(yamlFile, &fileCfg)
	if err != nil {
		// TODO add log

		return &config.Config{}
	}

	return fileCfg.toDomainConfig()
}

func (a *Adapter) ClearCurrentContext() error {
	err := os.Remove(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return nil
}

func (a *Adapter) StoreContext(cfg *config.Config) error {
	fileCfg := fileConfigFromDomainConfig(cfg)

	yamlBytes, err := yaml.Marshal(&fileCfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, yamlBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (a *Adapter) String(cfg *config.Config) (string, error) {
	fileCfg := fileConfigFromDomainConfig(cfg).fillUnsetSettings()

	str, err := fileCfg.string()
	if err != nil {
		return "", err
	}

	return str, nil
}

func (a *Adapter) StringJson(cfg *config.Config) (string, error) {
	fileCfg := fileConfigFromDomainConfig(cfg).fillUnsetSettings()

	str, err := fileCfg.stringJson()
	if err != nil {
		return "", err
	}

	return str, nil
}

func (a *Adapter) StringYaml(cfg *config.Config) (string, error) {
	fileCfg := fileConfigFromDomainConfig(cfg).fillUnsetSettings()

	str, err := fileCfg.stringYaml()
	if err != nil {
		return "", err
	}

	return str, nil
}
