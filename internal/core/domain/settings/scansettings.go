package settings

import (
	"reflect"
)

func (s *ScanSettings) HasBlackBoxSettings() bool {
	return s.BlackBoxEnabled
}

func (s *ScanSettings) IsEmpty() bool {
	return reflect.DeepEqual(s, &ScanSettings{})
}

type ScanSettings struct {
	ProjectName           string                `json:"projectName"`
	Languages             []string              `json:"languages"`
	WhiteBoxSettings      WhiteBoxSettings      `json:"whiteBoxSettings"`
	BlackBoxEnabled       bool                  `json:"blackBoxEnabled"`
	BlackBoxSettings      BlackBoxSettings      `json:"blackBoxSettings"`
	DotNetSettings        DotNetSettings        `json:"dotNetSettings"`
	GoSettings            GoSettings            `json:"goSettings"`
	JavaScriptSettings    JavaScriptSettings    `json:"javaScriptSettings"`
	JavaSettings          JavaSettings          `json:"javaSettings"`
	PhpSettings           PhpSettings           `json:"phpSettings"`
	PmTaintSettings       PmTaintSettings       `json:"pmTaintSettings"`
	PythonSettings        PythonSettings        `json:"pythonSettings"`
	RubySettings          RubySettings          `json:"rubySettings"`
	PygrepSettings        PygrepSettings        `json:"pygrepSettings"`
	ScaSettings           ScaSettings           `json:"scaSettings"`
	ReportAfterScan       ReportAfterScan       `json:"reportAfterScan"`
	SkipGitIgnoreFiles    bool                  `json:"skipGitIgnoreFiles"`
	VcsConnectionSettings VcsConnectionSettings `json:"vcsConnectionSettings"`
}

type WhiteBoxSettings struct {
	StaticCodeAnalysisEnabled            bool `json:"staticCodeAnalysisEnabled"`
	PatternMatchingEnabled               bool `json:"patternMatchingEnabled"`
	SearchForVulnerableComponentsEnabled bool `json:"searchForVulnerableComponentsEnabled"`
	SearchForConfigurationFlawsEnabled   bool `json:"searchForConfigurationFlawsEnabled"`
	SearchWithScaEnabled                 bool `json:"searchWithScaEnabled"`
	SecretDetectionEnabled               bool `json:"secretDetectionEnabled"`
	SearchForMaliciousCodeEnabled        bool `json:"searchForMaliciousCodeEnabled"`
}

type BlackBoxSettings struct {
	Site                  string                  `json:"site"`
	Level                 string                  `json:"level"`
	ScanScope             string                  `json:"scanScope"`
	SslCheck              bool                    `json:"sslCheck"`
	RunAutocheckAfterScan bool                    `json:"runAutocheckAfterScan"`
	AdditionalHttpHeaders []HTTPHeader            `json:"additionalHttpHeaders"`
	WhiteListedAddresses  []AddressEntry          `json:"whiteListedAddresses"`
	BlackListedAddresses  []AddressEntry          `json:"blackListedAddresses"`
	Authentication        *BlackBoxAuthentication `json:"authentication"`
	ProxySettings         *BlackBoxProxySettings  `json:"proxySettings"`
}

type HTTPHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AddressEntry struct {
	Address string `json:"address"`
	Format  string `json:"format"`
}

type BlackBoxAuthentication struct {
	Type   string              `json:"type"`
	Cookie *BlackBoxCookieAuth `json:"cookie"`
	Form   *BlackBoxFormAuth   `json:"form"`
	Http   *BlackBoxHTTPAuth   `json:"http"`
}

type BlackBoxCookieAuth struct {
	Cookie             string `json:"cookie"`
	ValidationAddress  string `json:"validationAddress"`
	ValidationTemplate string `json:"validationTemplate"`
}

type BlackBoxFormAuth struct {
	FormDetection      string `json:"formDetection"`
	FormAddress        string `json:"formAddress"`
	FormXPath          string `json:"formXPath"`
	Login              string `json:"login"`
	LoginKey           string `json:"loginKey"`
	Password           string `json:"password"`
	PasswordKey        string `json:"passwordKey"`
	ValidationTemplate string `json:"validationTemplate"`
}

type BlackBoxHTTPAuth struct {
	Login             string `json:"login"`
	Password          string `json:"password"`
	ValidationAddress string `json:"validationAddress"`
}

type BlackBoxProxySettings struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Type     string `json:"type"`
}

type DotNetSettings struct {
	ProjectType                           string `json:"projectType"`
	SolutionFile                          string `json:"solutionFile"`
	WebSiteFolder                         string `json:"webSiteFolder"`
	LaunchParameters                      string `json:"launchParameters"`
	UseAvailablePublicAndProtectedMethods bool   `json:"useAvailablePublicAndProtectedMethods"`
	DownloadDependencies                  bool   `json:"downloadDependencies"`
}

type GoSettings struct {
	LaunchParameters                      string `json:"launchParameters"`
	UseAvailablePublicAndProtectedMethods bool   `json:"useAvailablePublicAndProtectedMethods"`
}

type JavaScriptSettings struct {
	LaunchParameters                      string `json:"launchParameters"`
	UseAvailablePublicAndProtectedMethods bool   `json:"useAvailablePublicAndProtectedMethods"`
	DownloadDependencies                  bool   `json:"downloadDependencies"`
	UseTaintAnalysis                      bool   `json:"useTaintAnalysis"`
	UseJsaAnalysis                        bool   `json:"useJsaAnalysis"`
}

type JavaSettings struct {
	Parameters                            string `json:"parameters"`
	UnpackUserPackages                    bool   `json:"unpackUserPackages"`
	UserPackagePrefixes                   string `json:"userPackagePrefixes"`
	Version                               string `json:"version"`
	LaunchParameters                      string `json:"launchParameters"`
	UseAvailablePublicAndProtectedMethods bool   `json:"useAvailablePublicAndProtectedMethods"`
	DownloadDependencies                  bool   `json:"downloadDependencies"`
	DependenciesPath                      string `json:"dependenciesPath"`
}

type PhpSettings struct {
	LaunchParameters                      string `json:"launchParameters"`
	UseAvailablePublicAndProtectedMethods bool   `json:"useAvailablePublicAndProtectedMethods"`
	DownloadDependencies                  bool   `json:"downloadDependencies"`
}

type PmTaintSettings struct {
	LaunchParameters                      string `json:"launchParameters"`
	UseAvailablePublicAndProtectedMethods bool   `json:"useAvailablePublicAndProtectedMethods"`
}

type PythonSettings struct {
	LaunchParameters                      string `json:"launchParameters"`
	UseAvailablePublicAndProtectedMethods bool   `json:"useAvailablePublicAndProtectedMethods"`
	DownloadDependencies                  bool   `json:"downloadDependencies"`
	DependenciesPath                      string `json:"dependenciesPath"`
}

type RubySettings struct {
	LaunchParameters                      string `json:"launchParameters"`
	UseAvailablePublicAndProtectedMethods bool   `json:"useAvailablePublicAndProtectedMethods"`
}

type PygrepSettings struct {
	RulesDirPath     string `json:"rulesDirPath"`
	LaunchParameters string `json:"launchParameters"`
}

type ScaSettings struct {
	LaunchParameters       string `json:"launchParameters"`
	BuildDependenciesGraph bool   `json:"buildDependenciesGraph"`
}

type ReportAfterScan struct {
	Enabled         bool     `json:"enabled"`
	MailProfileId   string   `json:"mailProfileId"`
	EmailRecipients []string `json:"emailRecipients"`
}

type VcsConnectionSettings struct {
	RepositoryUrl                string `json:"repositoryUrl"`
	Login                        string `json:"login"`
	Password                     string `json:"password"`
	SshKey                       string `json:"sshKey"`
	IncludeSubmodules            bool   `json:"includeSubmodules"`
	SourceControlCredentialsType string `json:"sourceControlCredentialsType"`
	RepositoryType               string `json:"repositoryType"`
	CredentialsId                string `json:"credentialsId"`
	AzurePersonalAccessToken     string `json:"azurePersonalAccessToken"`
}
