package settings

// Language is an AI programming language group, spelled as the API spells it in licenses and project settings.
type Language string

const (
	LanguageNone          Language = "None"
	LanguageCAndCPlusPlus Language = "CAndCPlusPlus"
	LanguageCSharp        Language = "CSharp"
	LanguageCSharpWinOnly Language = "CSharpWinOnly"
	LanguageDart          Language = "Dart" // AIE 6.1+
	LanguageGo            Language = "Go"
	LanguageJava          Language = "Java"
	LanguageJavaScript    Language = "JavaScript"
	LanguageKotlin        Language = "Kotlin"
	LanguageObjectiveC    Language = "ObjectiveC"
	LanguageOneC          Language = "OneC" // AIE 6.0+
	LanguagePhp           Language = "Php"
	LanguagePython        Language = "Python"
	LanguageRuby          Language = "Ruby"
	LanguageScala         Language = "Scala"
	LanguageSolidity      Language = "Solidity"
	LanguageSql           Language = "Sql"
	LanguageSwift         Language = "Swift"
)
