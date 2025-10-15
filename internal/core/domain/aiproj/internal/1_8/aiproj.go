package aiproj1_8

// Generated from https://www.schemastore.org/aiproj-1.8.json

type AIProj struct {
	Schema                 *string                 `json:"$schema,omitempty"`
	Version                string                  `json:"Version"`
	BlackBoxSettings       *BlackBoxSettings       `json:"BlackBoxSettings,omitempty"`
	DotNetSettings         *DotNetSettings         `json:"DotNetSettings,omitempty"`
	GoSettings             *GoSettings             `json:"GoSettings,omitempty"`
	JavaSettings           *JavaSettings           `json:"JavaSettings,omitempty"`
	JavaScriptSettings     *JavaScriptSettings     `json:"JavaScriptSettings,omitempty"`
	PhpSettings            *PhpSettings            `json:"PhpSettings,omitempty"`
	PmTaintSettings        *PmTaintSettings        `json:"PmTaintSettings,omitempty"`
	PygrepSettings         *PygrepSettings         `json:"PygrepSettings,omitempty"`
	PythonSettings         *PythonSettings         `json:"PythonSettings,omitempty"`
	MailingProjectSettings *MailingProjectSettings `json:"MailingProjectSettings,omitempty"`
	RubySettings           *RubySettings           `json:"RubySettings,omitempty"`
	ScaSettings            *ScaSettings            `json:"ScaSettings,omitempty"`
	ProgrammingLanguages   []string                `json:"ProgrammingLanguages"`
	ProjectName            string                  `json:"ProjectName"`
	BranchName             *string                 `json:"BranchName,omitempty"`
	ScanModules            []string                `json:"ScanModules"`
	SkipGitIgnoreFiles     bool                    `json:"SkipGitIgnoreFiles,omitempty"`
	ApplyAllPMRules        bool                    `json:"ApplyAllPMRules,omitempty"`
	UseSecurityPolicies    bool                    `json:"UseSecurityPolicies,omitempty"`
}

type BlackBoxSettings struct {
	AdditionalHttpHeaders *[]HTTPHeader   `json:"AdditionalHttpHeaders,omitempty"`
	WhiteListedAddresses  *[]AddressEntry `json:"WhiteListedAddresses,omitempty"`
	BlackListedAddresses  *[]AddressEntry `json:"BlackListedAddresses,omitempty"`
	Authentication        *Authentication `json:"Authentication,omitempty"`
	Level                 *string         `json:"Level,omitempty"`
	ProxySettings         *ProxySettings  `json:"ProxySettings,omitempty"`
	RunAutocheckAfterScan *bool           `json:"RunAutocheckAfterScan,omitempty"`
	ScanScope             *string         `json:"ScanScope,omitempty"`
	Site                  *string         `json:"Site,omitempty"`
	SslCheck              *bool           `json:"SslCheck,omitempty"`
}

type HTTPHeader struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type AddressEntry struct {
	Address string `json:"Address"`
	Format  string `json:"Format"`
}

type Authentication struct {
	Type   string      `json:"Type"`
	Cookie *CookieAuth `json:"Cookie,omitempty"`
	Form   *FormAuth   `json:"Form,omitempty"`
	Http   *HttpAuth   `json:"Http,omitempty"`
}

type CookieAuth struct {
	Cookie             string `json:"Cookie"`
	ValidationAddress  string `json:"ValidationAddress"`
	ValidationTemplate string `json:"ValidationTemplate"`
}

type FormAuth struct {
	FormDetection      string  `json:"FormDetection"`
	FormAddress        string  `json:"FormAddress"`
	FormXPath          *string `json:"FormXPath,omitempty"`
	Login              string  `json:"Login"`
	LoginKey           *string `json:"LoginKey,omitempty"`
	Password           string  `json:"Password"`
	PasswordKey        *string `json:"PasswordKey,omitempty"`
	ValidationTemplate string  `json:"ValidationTemplate"`
}

type HttpAuth struct {
	Login             string `json:"Login"`
	Password          string `json:"Password"`
	ValidationAddress string `json:"ValidationAddress"`
}

type ProxySettings struct {
	Enabled  bool    `json:"Enabled"`
	Host     *string `json:"Host,omitempty"`
	Login    *string `json:"Login,omitempty"`
	Password *string `json:"Password,omitempty"`
	Port     int     `json:"Port"`
	Type     string  `json:"Type"`
}

type DotNetSettings struct {
	ProjectType             *string `json:"ProjectType,omitempty"`
	SolutionFile            *string `json:"SolutionFile,omitempty"`
	UsePublicAnalysisMethod *bool   `json:"UsePublicAnalysisMethod,omitempty"`
	DownloadDependencies    *bool   `json:"DownloadDependencies,omitempty"`
	CustomParameters        *string `json:"CustomParameters,omitempty"`
	DslRulesRelativePath    *string `json:"DslRulesRelativePath,omitempty"`
}

type GoSettings struct {
	UsePublicAnalysisMethod *bool   `json:"UsePublicAnalysisMethod,omitempty"`
	CustomParameters        *string `json:"CustomParameters,omitempty"`
	DslRulesRelativePath    *string `json:"DslRulesRelativePath,omitempty"`
}

type JavaSettings struct {
	Parameters              *string `json:"Parameters,omitempty"`
	UnpackUserPackages      bool    `json:"UnpackUserPackages"`
	UserPackagePrefixes     *string `json:"UserPackagePrefixes,omitempty"`
	Version                 string  `json:"Version"`
	UsePublicAnalysisMethod *bool   `json:"UsePublicAnalysisMethod,omitempty"`
	DownloadDependencies    *bool   `json:"DownloadDependencies,omitempty"`
	DependenciesPath        *string `json:"DependenciesPath,omitempty"`
	CustomParameters        *string `json:"CustomParameters,omitempty"`
	DslRulesRelativePath    *string `json:"DslRulesRelativePath,omitempty"`
}

type JavaScriptSettings struct {
	UsePublicAnalysisMethod *bool   `json:"UsePublicAnalysisMethod,omitempty"`
	UseTaintAnalysis        *bool   `json:"UseTaintAnalysis,omitempty"`
	UseJsaAnalysis          *bool   `json:"UseJsaAnalysis,omitempty"`
	DownloadDependencies    *bool   `json:"DownloadDependencies,omitempty"`
	CustomParameters        *string `json:"CustomParameters,omitempty"`
	DslRulesRelativePath    *string `json:"DslRulesRelativePath,omitempty"`
}

type PhpSettings struct {
	UsePublicAnalysisMethod *bool   `json:"UsePublicAnalysisMethod,omitempty"`
	DownloadDependencies    *bool   `json:"DownloadDependencies,omitempty"`
	CustomParameters        *string `json:"CustomParameters,omitempty"`
	DslRulesRelativePath    *string `json:"DslRulesRelativePath,omitempty"`
}

type PmTaintSettings struct {
	UsePublicAnalysisMethod *bool    `json:"UsePublicAnalysisMethod,omitempty"`
	PMGroups                []string `json:"PMGroups,omitempty"`
	CustomParameters        *string  `json:"CustomParameters,omitempty"`
}

type PygrepSettings struct {
	CustomParameters *string `json:"CustomParameters,omitempty"`
	RulesDirPath     *string `json:"RulesDirPath,omitempty"`
}

type PythonSettings struct {
	UsePublicAnalysisMethod *bool   `json:"UsePublicAnalysisMethod,omitempty"`
	DownloadDependencies    *bool   `json:"DownloadDependencies,omitempty"`
	DependenciesPath        *string `json:"DependenciesPath,omitempty"`
	CustomParameters        *string `json:"CustomParameters,omitempty"`
	DslRulesRelativePath    *string `json:"DslRulesRelativePath,omitempty"`
}

type MailingProjectSettings struct {
	Enabled         bool     `json:"Enabled"`
	MailProfileName *string  `json:"MailProfileName,omitempty"`
	EmailRecipients []string `json:"EmailRecipients,omitempty"`
}

type RubySettings struct {
	UsePublicAnalysisMethod *bool   `json:"UsePublicAnalysisMethod,omitempty"`
	CustomParameters        *string `json:"CustomParameters,omitempty"`
	DslRulesRelativePath    *string `json:"DslRulesRelativePath,omitempty"`
}

type ScaSettings struct {
	CustomParameters       *string `json:"CustomParameters,omitempty"`
	BuildDependenciesGraph *bool   `json:"BuildDependenciesGraph,omitempty"`
}
