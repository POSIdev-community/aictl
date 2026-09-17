## aictl get scan report with-filters markdown

Get scan report in Markdown format

### Synopsis

Download the scan report in Markdown format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.

```
aictl get scan report with-filters markdown <scan-id> [flags]
```

### Examples

```
  aictl get scan report with-filters markdown <scan-id> -o ./out.markdown --level-high
  aictl get scan report with-filters markdown <scan-id> -o ./out.markdown -f --level-high
```

### Options

```
  -h, --help   help for markdown
```

### Options inherited from parent commands

```
      --cacert string             Path to PEM file with CA certificate(s) to trust (appended to system roots)
      --conditional               Include conditional issues
  -V, --debug                     Debug output (includes error chains)
  -f, --force                     Overwrite existing output file
      --found-prev-scan           Include issues found in a previous scan
      --found-this-scan           Include issues found in this scan
      --include-comments          Include comments in the report
      --include-dfd               Include data flow diagrams in the report
      --include-glossary          Include glossary in the report
      --language stringArray      Language filter (repeatable): CAndCPlusPlus, CSharp, Dart, Go, Java, JavaScript, Kotlin, ObjectiveC, OneC, Php, Python, Ruby, Scala, Solidity, Sql, Swift
      --level-high                Include high severity issues
      --level-low                 Include low severity issues
      --level-medium              Include medium severity issues
      --level-potential           Include potential severity issues
      --localization string       Report localization language: 'en' or 'ru' (default "en")
  -l, --log-path string           Log file path
      --mode-entry-point          Include entry-point scan mode issues
      --mode-others               Include other scan mode issues
      --mode-public-methods       Include public-methods scan mode issues
      --mode-root-function        Include root-function scan mode issues
      --non-conditional           Include non-conditional issues
      --non-suppressed            Include non-suppressed issues
      --only-favorite             Include only favorite issues
  -o, --output string             Output file path
  -p, --project-id string         Project id (overrides context)
      --scan-module stringArray   Scan module filter (repeatable): StaticCodeAnalysis, PatternMatching, Components, SoftwareCompositionAnalysis, Configuration, MaliciousCodeDetection, SecretDetection, BlackBox
      --second-level              Include second-level issues
      --status-confirmed          Include confirmed issues
      --status-confirmed-auto     Include auto-confirmed issues
      --status-rejected           Include rejected issues
      --status-undefined          Include issues with undefined confirmation status
      --suppressed                Include suppressed issues
      --suspected                 Include suspected issues
      --tls-skip                  Skip TLS certificate verification
  -t, --token string              AI server access token (overrides context)
      --type stringArray          Vulnerability type filter (repeatable)
  -u, --uri string                AI server URI (overrides context)
  -v, --verbose                   Verbose output (operational details)
```

### SEE ALSO

* [aictl get scan report with-filters](aictl_get_scan_report_with-filters.md)	 - Get scan report with vulnerability filters

