package aiproj

import (
	v110 "github.com/POSIdev-community/aiproj/model/v1_10"
	v111 "github.com/POSIdev-community/aiproj/model/v1_11"
	v19 "github.com/POSIdev-community/aiproj/model/v1_9"
)

type V19 = v19.AIProj
type V110 = v110.AIProj
type V111 = v111.AIProj

type ScanModule string

const (
	ScanModuleConfiguration               ScanModule = "Configuration"
	ScanModuleComponents                  ScanModule = "Components"
	ScanModuleBlackBox                    ScanModule = "BlackBox"
	ScanModulePatternMatching             ScanModule = "PatternMatching"
	ScanModuleStaticCodeAnalysis          ScanModule = "StaticCodeAnalysis"
	ScanModuleSoftwareCompositionAnalysis ScanModule = "SoftwareCompositionAnalysis"
	ScanModuleSecretDetection             ScanModule = "SecretDetection"
	ScanModuleMaliciousCodeDetection      ScanModule = "MaliciousCodeDetection"
)

type ProgrammingLanguage string

const (
	ProgrammingLanguageCSharpWindowsLinux ProgrammingLanguage = "CSharp (Windows, Linux)"
	ProgrammingLanguageCSharpWindows      ProgrammingLanguage = "CSharp (Windows)"
)

type Result struct {
	Version string
	V19     *V19
	V110    *V110
	V111    *V111
}

func (r *Result) AiprojVersion() string {
	if r == nil {
		return ""
	}

	return r.Version
}
