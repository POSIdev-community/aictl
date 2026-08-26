## aictl get scan report

Get scan report

### Synopsis

Download a scan report by custom template name or use a built-in format subcommand. For custom templates, pass report name and scan id. Output path via -o; use -f to overwrite.

```
aictl get scan report <report-name> <scan-id> [flags]
```

### Examples

```
  aictl get scan report MyTemplate <scan-id> -o ./out.html
  aictl get scan report sarif <scan-id> -o ./out.sarif
```

### Options

```
  -f, --force                 Overwrite existing output file
  -h, --help                  help for report
      --include-comments      Include comments in the report
      --include-dfd           Include data flow diagrams in the report
      --include-glossary      Include glossary in the report
      --localization string   Report localization language: 'en' or 'ru' (default "en")
  -o, --output string         Output file path
```

### Options inherited from parent commands

```
      --cacert string       Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -l, --log-path string     Log file path
  -p, --project-id string   Project id (overrides context)
      --tls-skip            Skip TLS certificate verification
  -t, --token string        AI server access token (overrides context)
  -u, --uri string          AI server URI (overrides context)
  -v, --verbose             Verbose output
```

### SEE ALSO

* [aictl get scan](aictl_get_scan.md)	 - Get scan
* [aictl get scan report autocheck](aictl_get_scan_report_autocheck.md)	 - Get scan report in AutoCheck format
* [aictl get scan report gitlab](aictl_get_scan_report_gitlab.md)	 - Get scan report in GitLab format
* [aictl get scan report json](aictl_get_scan_report_json.md)	 - Get scan report in JSON format
* [aictl get scan report json-v2](aictl_get_scan_report_json-v2.md)	 - Get scan report in JSON v2 format
* [aictl get scan report markdown](aictl_get_scan_report_markdown.md)	 - Get scan report in Markdown format
* [aictl get scan report nist](aictl_get_scan_report_nist.md)	 - Get scan report in NIST format
* [aictl get scan report oud4](aictl_get_scan_report_oud4.md)	 - Get scan report in OUD4 format
* [aictl get scan report owasp](aictl_get_scan_report_owasp.md)	 - Get scan report in OWASP format
* [aictl get scan report owaspm](aictl_get_scan_report_owaspm.md)	 - Get scan report in OWASP Maturity format
* [aictl get scan report pcidss](aictl_get_scan_report_pcidss.md)	 - Get scan report in PCI DSS format
* [aictl get scan report plain](aictl_get_scan_report_plain.md)	 - Get scan report in plain text format
* [aictl get scan report sans](aictl_get_scan_report_sans.md)	 - Get scan report in SANS format
* [aictl get scan report sarif](aictl_get_scan_report_sarif.md)	 - Get scan report in SARIF format
* [aictl get scan report xml](aictl_get_scan_report_xml.md)	 - Get scan report in XML format

