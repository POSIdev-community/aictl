package settings

import (
	"reflect"

	"github.com/POSIdev-community/aictl/internal/core/domain/aiproj"
)

func (s *ScanSettings) UpdateFromAIProj(aiproj *aiproj.AIProj) *ScanSettings {
	if aiproj.ProjectName != "" {
		s.ProjectName = aiproj.ProjectName
	}

	if aiproj.ProgrammingLanguages != nil {
		s.Languages = make([]string, len(aiproj.ProgrammingLanguages))
		for i, lang := range aiproj.ProgrammingLanguages {
			switch lang {
			case "CSharp (Windows, Linux)":
				s.Languages[i] = "CSharp"
			case "CSharp (Windows)":
				s.Languages[i] = "CSharpWinOnly"
			default:
				s.Languages[i] = lang
			}
		}
	}

	s.SkipGitIgnoreFiles = aiproj.SkipGitIgnoreFiles

	// Override DotNetSettings if present in AIProj
	if aiproj.DotNetSettings != nil {
		if aiproj.DotNetSettings.ProjectType != nil {
			s.DotNetSettings.ProjectType = *aiproj.DotNetSettings.ProjectType
		}
		if aiproj.DotNetSettings.SolutionFile != nil {
			s.DotNetSettings.SolutionFile = *aiproj.DotNetSettings.SolutionFile
		}
		if aiproj.DotNetSettings.UsePublicAnalysisMethod != nil {
			s.DotNetSettings.UseAvailablePublicAndProtectedMethods = *aiproj.DotNetSettings.UsePublicAnalysisMethod
		}
		if aiproj.DotNetSettings.DownloadDependencies != nil {
			s.DotNetSettings.DownloadDependencies = *aiproj.DotNetSettings.DownloadDependencies
		}
		if aiproj.DotNetSettings.CustomParameters != nil {
			s.DotNetSettings.LaunchParameters = *aiproj.DotNetSettings.CustomParameters
		}
	}

	// Override GoSettings if present in AIProj
	if aiproj.GoSettings != nil {
		if aiproj.GoSettings.CustomParameters != nil {
			s.GoSettings.LaunchParameters = *aiproj.GoSettings.CustomParameters
		}
		if aiproj.GoSettings.UsePublicAnalysisMethod != nil {
			s.GoSettings.UseAvailablePublicAndProtectedMethods = *aiproj.GoSettings.UsePublicAnalysisMethod
		}
	}

	// Override JavaScriptSettings if present in AIProj
	if aiproj.JavaScriptSettings != nil {
		if aiproj.JavaScriptSettings.CustomParameters != nil {
			s.JavaScriptSettings.LaunchParameters = *aiproj.JavaScriptSettings.CustomParameters
		}
		if aiproj.JavaScriptSettings.UsePublicAnalysisMethod != nil {
			s.JavaScriptSettings.UseAvailablePublicAndProtectedMethods = *aiproj.JavaScriptSettings.UsePublicAnalysisMethod
		}
		if aiproj.JavaScriptSettings.DownloadDependencies != nil {
			s.JavaScriptSettings.DownloadDependencies = *aiproj.JavaScriptSettings.DownloadDependencies
		}
		if aiproj.JavaScriptSettings.UseTaintAnalysis != nil {
			s.JavaScriptSettings.UseTaintAnalysis = *aiproj.JavaScriptSettings.UseTaintAnalysis
		}
		if aiproj.JavaScriptSettings.UseJsaAnalysis != nil {
			s.JavaScriptSettings.UseJsaAnalysis = *aiproj.JavaScriptSettings.UseJsaAnalysis
		}
	}

	// Override JavaSettings if present in AIProj
	if aiproj.JavaSettings != nil {
		if aiproj.JavaSettings.Parameters != nil {
			s.JavaSettings.Parameters = *aiproj.JavaSettings.Parameters
		}

		s.JavaSettings.UnpackUserPackages = aiproj.JavaSettings.UnpackUserPackages

		if aiproj.JavaSettings.UserPackagePrefixes != nil {
			s.JavaSettings.UserPackagePrefixes = *aiproj.JavaSettings.UserPackagePrefixes
		}

		if aiproj.JavaSettings.Version == "" {
			s.JavaSettings.Version = "21"
		}

		s.JavaSettings.Version = "v1_" + aiproj.JavaSettings.Version

		if aiproj.JavaSettings.CustomParameters != nil {
			s.JavaSettings.LaunchParameters = *aiproj.JavaSettings.CustomParameters
		}
		if aiproj.JavaSettings.UsePublicAnalysisMethod != nil {
			s.JavaSettings.UseAvailablePublicAndProtectedMethods = *aiproj.JavaSettings.UsePublicAnalysisMethod
		}
		if aiproj.JavaSettings.DownloadDependencies != nil {
			s.JavaSettings.DownloadDependencies = *aiproj.JavaSettings.DownloadDependencies
		}
		if aiproj.JavaSettings.DependenciesPath != nil {
			s.JavaSettings.DependenciesPath = *aiproj.JavaSettings.DependenciesPath
		}
	}

	// Override PhpSettings if present in AIProj
	if aiproj.PhpSettings != nil {
		if aiproj.PhpSettings.CustomParameters != nil {
			s.PhpSettings.LaunchParameters = *aiproj.PhpSettings.CustomParameters
		}
		if aiproj.PhpSettings.UsePublicAnalysisMethod != nil {
			s.PhpSettings.UseAvailablePublicAndProtectedMethods = *aiproj.PhpSettings.UsePublicAnalysisMethod
		}
		if aiproj.PhpSettings.DownloadDependencies != nil {
			s.PhpSettings.DownloadDependencies = *aiproj.PhpSettings.DownloadDependencies
		}
	}

	// Override PmTaintSettings if present in AIProj
	if aiproj.PmTaintSettings != nil {
		if aiproj.PmTaintSettings.CustomParameters != nil {
			s.PmTaintSettings.LaunchParameters = *aiproj.PmTaintSettings.CustomParameters
		}
		if aiproj.PmTaintSettings.UsePublicAnalysisMethod != nil {
			s.PmTaintSettings.UseAvailablePublicAndProtectedMethods = *aiproj.PmTaintSettings.UsePublicAnalysisMethod
		}
	}

	// Override PythonSettings if present in AIProj
	if aiproj.PythonSettings != nil {
		if aiproj.PythonSettings.CustomParameters != nil {
			s.PythonSettings.LaunchParameters = *aiproj.PythonSettings.CustomParameters
		}
		if aiproj.PythonSettings.UsePublicAnalysisMethod != nil {
			s.PythonSettings.UseAvailablePublicAndProtectedMethods = *aiproj.PythonSettings.UsePublicAnalysisMethod
		}
		if aiproj.PythonSettings.DownloadDependencies != nil {
			s.PythonSettings.DownloadDependencies = *aiproj.PythonSettings.DownloadDependencies
		}
		if aiproj.PythonSettings.DependenciesPath != nil {
			s.PythonSettings.DependenciesPath = *aiproj.PythonSettings.DependenciesPath
		}
	}

	// Override RubySettings if present in AIProj
	if aiproj.RubySettings != nil {
		if aiproj.RubySettings.CustomParameters != nil {
			s.RubySettings.LaunchParameters = *aiproj.RubySettings.CustomParameters
		}
		if aiproj.RubySettings.UsePublicAnalysisMethod != nil {
			s.RubySettings.UseAvailablePublicAndProtectedMethods = *aiproj.RubySettings.UsePublicAnalysisMethod
		}
	}

	// Override PygrepSettings if present in AIProj
	if aiproj.PygrepSettings != nil {
		if aiproj.PygrepSettings.RulesDirPath != nil {
			s.PygrepSettings.RulesDirPath = *aiproj.PygrepSettings.RulesDirPath
		}
		if aiproj.PygrepSettings.CustomParameters != nil {
			s.PygrepSettings.LaunchParameters = *aiproj.PygrepSettings.CustomParameters
		}
	}

	// Override ScaSettings if present in AIProj
	if aiproj.ScaSettings != nil {
		if aiproj.ScaSettings.CustomParameters != nil {
			s.ScaSettings.LaunchParameters = *aiproj.ScaSettings.CustomParameters
		}
		if aiproj.ScaSettings.BuildDependenciesGraph != nil {
			s.ScaSettings.BuildDependenciesGraph = *aiproj.ScaSettings.BuildDependenciesGraph
		}
	}

	// Override ReportAfterScan if present in AIProj
	if aiproj.MailingProjectSettings != nil {
		s.ReportAfterScan.Enabled = aiproj.MailingProjectSettings.Enabled
		if aiproj.MailingProjectSettings.MailProfileName != nil {
			s.ReportAfterScan.MailProfileId = *aiproj.MailingProjectSettings.MailProfileName
		}
		if len(aiproj.MailingProjectSettings.EmailRecipients) > 0 {
			s.ReportAfterScan.EmailRecipients = aiproj.MailingProjectSettings.EmailRecipients
		}
	}

	return s
}

func (s *ScanSettings) IsEmpty() bool {
	return reflect.DeepEqual(s, &ScanSettings{})
}

type ScanSettings struct {
	ProjectName           string                `json:"projectName"`
	Languages             []string              `json:"languages"`
	WhiteBoxSettings      WhiteBoxSettings      `json:"whiteBoxSettings"`
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
