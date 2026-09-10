## aictl get scan report with-filters

Get scan report with vulnerability filters

### Synopsis

Download a scan report with vulnerability filters applied (useFilters=true).
Requires at least one filter flag. Same format subcommands and output flags as get scan report.
Boolean filter flags send true only when present; unset flags are omitted from the API payload.
Array filters (--type, --language, --scan-module) that are not passed are sent as empty arrays.

```
aictl get scan report with-filters <report-name> <scan-id> [flags]
```

### Examples

```
  aictl get scan report with-filters sarif <scan-id> -o ./out.sarif --level-high --level-medium
  aictl get scan report with-filters MyTemplate <scan-id> -o ./out.html --status-confirmed --language Java
```

### Options

```
      --conditional               Include conditional issues
      --found-prev-scan           Include issues found in a previous scan
      --found-this-scan           Include issues found in this scan
  -h, --help                      help for with-filters
      --language stringArray      Language filter (repeatable): CAndCPlusPlus, CSharp, Dart, Go, Java, JavaScript, Kotlin, ObjectiveC, OneC, Php, Python, Ruby, Scala, Solidity, Sql, Swift
      --level-high                Include high severity issues
      --level-low                 Include low severity issues
      --level-medium              Include medium severity issues
      --level-potential           Include potential severity issues
      --mode-entry-point          Include entry-point scan mode issues
      --mode-others               Include other scan mode issues
      --mode-public-methods       Include public-methods scan mode issues
      --mode-root-function        Include root-function scan mode issues
      --non-conditional           Include non-conditional issues
      --non-suppressed            Include non-suppressed issues
      --only-favorite             Include only favorite issues
      --scan-module stringArray   Scan module filter (repeatable): StaticCodeAnalysis, PatternMatching, Components, SoftwareCompositionAnalysis, Configuration, MaliciousCodeDetection, SecretDetection, BlackBox
      --second-level              Include second-level issues
      --status-confirmed          Include confirmed issues
      --status-confirmed-auto     Include auto-confirmed issues
      --status-rejected           Include rejected issues
      --status-undefined          Include issues with undefined confirmation status
      --suppressed                Include suppressed issues
      --suspected                 Include suspected issues
      --type stringArray          Vulnerability type filter (repeatable)
```

### Options inherited from parent commands

```
      --cacert string         Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -V, --debug                 Debug output (includes error chains)
  -f, --force                 Overwrite existing output file
      --include-comments      Include comments in the report
      --include-dfd           Include data flow diagrams in the report
      --include-glossary      Include glossary in the report
      --localization string   Report localization language: 'en' or 'ru' (default "en")
  -l, --log-path string       Log file path
  -o, --output string         Output file path
  -p, --project-id string     Project id (overrides context)
      --tls-skip              Skip TLS certificate verification
  -t, --token string          AI server access token (overrides context)
  -u, --uri string            AI server URI (overrides context)
  -v, --verbose               Verbose output (operational details)
```

### SEE ALSO

* [aictl get scan report](aictl_get_scan_report.md)	 - Get scan report
* [aictl get scan report with-filters autocheck](aictl_get_scan_report_with-filters_autocheck.md)	 - Get scan report in Autocheck format
* [aictl get scan report with-filters gitlab](aictl_get_scan_report_with-filters_gitlab.md)	 - Get scan report in GitLab format
* [aictl get scan report with-filters json](aictl_get_scan_report_with-filters_json.md)	 - Get scan report in JSON format
* [aictl get scan report with-filters json-v2](aictl_get_scan_report_with-filters_json-v2.md)	 - Get scan report in JSON v2 format
* [aictl get scan report with-filters markdown](aictl_get_scan_report_with-filters_markdown.md)	 - Get scan report in Markdown format
* [aictl get scan report with-filters nist](aictl_get_scan_report_with-filters_nist.md)	 - Get scan report in NIST format
* [aictl get scan report with-filters oud4](aictl_get_scan_report_with-filters_oud4.md)	 - Get scan report in OUD4 format
* [aictl get scan report with-filters owasp](aictl_get_scan_report_with-filters_owasp.md)	 - Get scan report in OWASP format
* [aictl get scan report with-filters owaspm](aictl_get_scan_report_with-filters_owaspm.md)	 - Get scan report in OWASP Mobile format
* [aictl get scan report with-filters pcidss](aictl_get_scan_report_with-filters_pcidss.md)	 - Get scan report in PCI DSS format
* [aictl get scan report with-filters plain](aictl_get_scan_report_with-filters_plain.md)	 - Get scan report in plain HTML format
* [aictl get scan report with-filters sans](aictl_get_scan_report_with-filters_sans.md)	 - Get scan report in SANS format
* [aictl get scan report with-filters sarif](aictl_get_scan_report_with-filters_sarif.md)	 - Get scan report in SARIF format
* [aictl get scan report with-filters xml](aictl_get_scan_report_with-filters_xml.md)	 - Get scan report in XML format

