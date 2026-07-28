## aictl get scan report json-v2

Get scan report in JSON v2 format

### Synopsis

Download the scan report in JSON v2 format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.

```
aictl get scan report json-v2 <scan-id> [flags]
```

### Examples

```
  aictl get scan report json-v2 <scan-id> -o ./out.json
  aictl get scan report json-v2 <scan-id> -o ./out.json -f
```

### Options

```
  -h, --help   help for json-v2
```

### Options inherited from parent commands

```
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
  -v, --verbose               Verbose output
```

### SEE ALSO

* [aictl get scan report](aictl_get_scan_report.md)	 - Get scan report

