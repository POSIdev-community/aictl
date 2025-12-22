package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

const (
	appDir     = "aictl"
	configFile = "/context.yaml"
)

type Adapter struct {
}

func NewContextAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) GetContextFromAictlFolder() *config.Config {
	configPath, err := getConfigPath()
	if err != nil {
		return &config.Config{}
	}

	if !fshelper.PathExists(configPath) {
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
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.Remove(configPath)
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

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(configDir, appDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directories: %w", err)
	}
	configPath := filepath.Join(dir, configFile)
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

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(configDir, appDir, configFile)

	return configPath, nil
}
