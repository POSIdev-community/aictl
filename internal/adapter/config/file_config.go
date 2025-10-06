package config

import (
	"encoding/json"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"strings"
)

type fileConfig struct {
	Uri       string `yaml:"uri" json:"uri"`
	Token     string `yaml:"token" json:"token"`
	TLSSkip   bool   `yaml:"tlsSkip" json:"tlsSkip"`
	ProjectId string `yaml:"projectId" json:"projectId"`
	BranchId  string `yaml:"branchId" json:"branchId"`
}

func fileConfigFromDomainConfig(config *config.Config) fileConfig {
	var projectId string
	if config.ProjectId() == uuid.Nil {
		projectId = ""
	} else {
		projectId = config.ProjectId().String()
	}

	var branchId string
	if config.BranchId() == uuid.Nil {
		branchId = ""
	} else {
		branchId = config.BranchId().String()
	}

	return fileConfig{
		Uri:       config.UriString(),
		Token:     config.Token(),
		TLSSkip:   config.TLSSkip(),
		ProjectId: projectId,
		BranchId:  branchId,
	}
}

func (fileCfg fileConfig) toDomainConfig() *config.Config {
	uri, _ := config.NewUri(fileCfg.Uri)

	projectId, err := uuid.Parse(fileCfg.ProjectId)
	if err != nil {
		projectId = uuid.Nil
	}

	branchId, err := uuid.Parse(fileCfg.BranchId)
	if err != nil {
		branchId = uuid.Nil
	}

	return config.NewConfig(uri, fileCfg.Token, fileCfg.TLSSkip, projectId, branchId)
}

func (fileCfg fileConfig) stringYaml() (string, error) {
	yamlBytes, err := yaml.Marshal(fileCfg)
	if err != nil {
		return "", err
	}

	yamlString := string(yamlBytes)

	return yamlString, nil
}

func (fileCfg fileConfig) string() (string, error) {
	yamlString, err := fileCfg.stringYaml()
	if err != nil {
		return "", err
	}

	yamlString = strings.Replace(yamlString, "tlsSkip", "tls-skip", 1)

	return yamlString, nil
}

func (fileCfg fileConfig) stringJson() (string, error) {
	jsonBytes, err := json.MarshalIndent(fileCfg, "", "    ")
	if err != nil {
		return "", err
	}

	jsonString := string(jsonBytes)

	return jsonString, nil
}

func (fileCfg fileConfig) fillUnsetSettings() fileConfig {
	if fileCfg.Uri == "" {
		fileCfg.Uri = "<unset>"
	}

	if fileCfg.Token == "" {
		fileCfg.Token = "<unset>"
	}

	if fileCfg.ProjectId == "" {
		fileCfg.ProjectId = "<unset>"
	}

	if fileCfg.BranchId == "" {
		fileCfg.BranchId = "<unset>"
	}

	return fileCfg
}
