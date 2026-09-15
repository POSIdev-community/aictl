package config

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

const (
	unsetValue  = "<unset>"
	maskedToken = "<masked>"
)

type fileConfig struct {
	Uri       string `yaml:"uri" json:"uri"`
	Token     string `yaml:"token" json:"token"`
	TLSSkip   bool   `yaml:"tlsSkip" json:"tlsSkip"`
	CACert    string `yaml:"caCert,omitempty" json:"caCert,omitempty"`
	ProjectId string `yaml:"projectId" json:"projectId"`
	BranchId  string `yaml:"branchId" json:"branchId"`
}

// fileConfigExport is used for machine-readable show (--json/--yaml):
// raw token is preserved; unset string fields become JSON/YAML null instead of "<unset>".
type fileConfigExport struct {
	Uri       *string `yaml:"uri" json:"uri"`
	Token     *string `yaml:"token" json:"token"`
	TLSSkip   bool    `yaml:"tlsSkip" json:"tlsSkip"`
	CACert    *string `yaml:"caCert" json:"caCert"`
	ProjectId *string `yaml:"projectId" json:"projectId"`
	BranchId  *string `yaml:"branchId" json:"branchId"`
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
		CACert:    config.CACertPath(),
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

	cfg := config.NewConfig(uri, fileCfg.Token, fileCfg.TLSSkip, projectId, branchId)
	_ = cfg.ApplyCACertPath(fileCfg.CACert)

	return cfg
}

func (fileCfg fileConfig) stringYaml() (string, error) {
	yamlBytes, err := yaml.Marshal(fileCfg.toExport())
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}

func (fileCfg fileConfig) string() (string, error) {
	display := fileCfg.maskToken().fillUnsetSettings()

	yamlBytes, err := yaml.Marshal(display)
	if err != nil {
		return "", err
	}

	yamlString := string(yamlBytes)
	yamlString = strings.Replace(yamlString, "tlsSkip", "tls-skip", 1)
	yamlString = strings.Replace(yamlString, "caCert", "ca-cert", 1)

	return yamlString, nil
}

func (fileCfg fileConfig) stringJson() (string, error) {
	jsonBytes, err := json.MarshalIndent(fileCfg.toExport(), "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (fileCfg fileConfig) fillUnsetSettings() fileConfig {
	if fileCfg.Uri == "" {
		fileCfg.Uri = unsetValue
	}

	if fileCfg.Token == "" {
		fileCfg.Token = unsetValue
	}

	if fileCfg.ProjectId == "" {
		fileCfg.ProjectId = unsetValue
	}

	if fileCfg.BranchId == "" {
		fileCfg.BranchId = unsetValue
	}

	if fileCfg.CACert == "" {
		fileCfg.CACert = unsetValue
	}

	return fileCfg
}

func (fileCfg fileConfig) maskToken() fileConfig {
	if fileCfg.Token != "" {
		fileCfg.Token = maskedToken
	}

	return fileCfg
}

func (fileCfg fileConfig) toExport() fileConfigExport {
	return fileConfigExport{
		Uri:       optionalString(fileCfg.Uri),
		Token:     optionalString(fileCfg.Token),
		TLSSkip:   fileCfg.TLSSkip,
		CACert:    optionalString(fileCfg.CACert),
		ProjectId: optionalString(fileCfg.ProjectId),
		BranchId:  optionalString(fileCfg.BranchId),
	}
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
