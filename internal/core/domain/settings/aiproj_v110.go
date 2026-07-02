package settings

import (
	v110 "github.com/POSIdev-community/aiproj/model/v1_10"
)

func (s *ScanSettings) UpdateFromV110(p *v110.AIProj) error {
	if len(p.ScanModules) > 0 {
		modules := make([]string, len(p.ScanModules))
		for i, module := range p.ScanModules {
			modules[i] = string(module)
		}

		if err := applyScanModules(s, modules); err != nil {
			return err
		}
	}

	if p.ProjectName != "" {
		s.ProjectName = p.ProjectName
	}

	if len(p.ProgrammingLanguages) > 0 {
		languages := make([]string, len(p.ProgrammingLanguages))
		for i, lang := range p.ProgrammingLanguages {
			languages[i] = string(lang)
		}
		applyLanguages(s, languages)
	}

	if p.SkipGitIgnoreFiles != nil {
		s.SkipGitIgnoreFiles = *p.SkipGitIgnoreFiles
	}

	if dotNet := p.DotNetSettings; dotNet != nil {
		if dotNet.ProjectType != nil {
			s.DotNetSettings.ProjectType = string(*dotNet.ProjectType)
		}
		if dotNet.SolutionFile != nil {
			s.DotNetSettings.SolutionFile = *dotNet.SolutionFile
		}
		if dotNet.UsePublicAnalysisMethod != nil {
			s.DotNetSettings.UseAvailablePublicAndProtectedMethods = *dotNet.UsePublicAnalysisMethod
		}
		if dotNet.DownloadDependencies != nil {
			s.DotNetSettings.DownloadDependencies = *dotNet.DownloadDependencies
		}
		if dotNet.CustomParameters != nil {
			s.DotNetSettings.LaunchParameters = *dotNet.CustomParameters
		}
	}

	if goSettings := p.GoSettings; goSettings != nil {
		if goSettings.CustomParameters != nil {
			s.GoSettings.LaunchParameters = *goSettings.CustomParameters
		}
		if goSettings.UsePublicAnalysisMethod != nil {
			s.GoSettings.UseAvailablePublicAndProtectedMethods = *goSettings.UsePublicAnalysisMethod
		}
	}

	if jsSettings := p.JavaScriptSettings; jsSettings != nil {
		if jsSettings.CustomParameters != nil {
			s.JavaScriptSettings.LaunchParameters = *jsSettings.CustomParameters
		}
		if jsSettings.UsePublicAnalysisMethod != nil {
			s.JavaScriptSettings.UseAvailablePublicAndProtectedMethods = *jsSettings.UsePublicAnalysisMethod
		}
		if jsSettings.DownloadDependencies != nil {
			s.JavaScriptSettings.DownloadDependencies = *jsSettings.DownloadDependencies
		}
		if jsSettings.UseTaintAnalysis != nil {
			s.JavaScriptSettings.UseTaintAnalysis = *jsSettings.UseTaintAnalysis
		}
		if jsSettings.UseJsaAnalysis != nil {
			s.JavaScriptSettings.UseJsaAnalysis = *jsSettings.UseJsaAnalysis
		}
	}

	if javaSettings := p.JavaSettings; javaSettings != nil {
		if javaSettings.Parameters != nil {
			s.JavaSettings.Parameters = *javaSettings.Parameters
		}
		if javaSettings.UnpackUserPackages != nil {
			s.JavaSettings.UnpackUserPackages = *javaSettings.UnpackUserPackages
		}
		if javaSettings.UserPackagePrefixes != nil {
			s.JavaSettings.UserPackagePrefixes = *javaSettings.UserPackagePrefixes
		}
		if javaSettings.Version != nil {
			s.JavaSettings.Version = "v1_" + string(*javaSettings.Version)
		}
		if javaSettings.CustomParameters != nil {
			s.JavaSettings.LaunchParameters = *javaSettings.CustomParameters
		}
		if javaSettings.UsePublicAnalysisMethod != nil {
			s.JavaSettings.UseAvailablePublicAndProtectedMethods = *javaSettings.UsePublicAnalysisMethod
		}
		if javaSettings.DownloadDependencies != nil {
			s.JavaSettings.DownloadDependencies = *javaSettings.DownloadDependencies
		}
		if javaSettings.DependenciesPath != nil {
			s.JavaSettings.DependenciesPath = *javaSettings.DependenciesPath
		}
	}

	if phpSettings := p.PhpSettings; phpSettings != nil {
		if phpSettings.CustomParameters != nil {
			s.PhpSettings.LaunchParameters = *phpSettings.CustomParameters
		}
		if phpSettings.UsePublicAnalysisMethod != nil {
			s.PhpSettings.UseAvailablePublicAndProtectedMethods = *phpSettings.UsePublicAnalysisMethod
		}
		if phpSettings.DownloadDependencies != nil {
			s.PhpSettings.DownloadDependencies = *phpSettings.DownloadDependencies
		}
	}

	if pmTaintSettings := p.PmTaintSettings; pmTaintSettings != nil {
		if pmTaintSettings.CustomParameters != nil {
			s.PmTaintSettings.LaunchParameters = *pmTaintSettings.CustomParameters
		}
		if pmTaintSettings.UsePublicAnalysisMethod != nil {
			s.PmTaintSettings.UseAvailablePublicAndProtectedMethods = *pmTaintSettings.UsePublicAnalysisMethod
		}
	}

	if pythonSettings := p.PythonSettings; pythonSettings != nil {
		if pythonSettings.CustomParameters != nil {
			s.PythonSettings.LaunchParameters = *pythonSettings.CustomParameters
		}
		if pythonSettings.UsePublicAnalysisMethod != nil {
			s.PythonSettings.UseAvailablePublicAndProtectedMethods = *pythonSettings.UsePublicAnalysisMethod
		}
		if pythonSettings.DownloadDependencies != nil {
			s.PythonSettings.DownloadDependencies = *pythonSettings.DownloadDependencies
		}
		if pythonSettings.DependenciesPath != nil {
			s.PythonSettings.DependenciesPath = *pythonSettings.DependenciesPath
		}
	}

	if rubySettings := p.RubySettings; rubySettings != nil {
		if rubySettings.CustomParameters != nil {
			s.RubySettings.LaunchParameters = *rubySettings.CustomParameters
		}
		if rubySettings.UsePublicAnalysisMethod != nil {
			s.RubySettings.UseAvailablePublicAndProtectedMethods = *rubySettings.UsePublicAnalysisMethod
		}
	}

	if scaSettings := p.ScaSettings; scaSettings != nil {
		if scaSettings.CustomParameters != nil {
			s.ScaSettings.LaunchParameters = *scaSettings.CustomParameters
		}
		if scaSettings.BuildDependenciesGraph != nil {
			s.ScaSettings.BuildDependenciesGraph = *scaSettings.BuildDependenciesGraph
		}
	}

	if mailingSettings := p.MailingProjectSettings; mailingSettings != nil {
		if mailingSettings.Enabled != nil {
			s.ReportAfterScan.Enabled = *mailingSettings.Enabled
		}
		if mailingSettings.MailProfileName != nil {
			s.ReportAfterScan.MailProfileId = *mailingSettings.MailProfileName
		}
		if len(mailingSettings.EmailRecipients) > 0 {
			s.ReportAfterScan.EmailRecipients = mailingSettings.EmailRecipients
		}
	}

	s.applyBlackBoxFromV110(p.BlackBoxSettings)

	return nil
}
