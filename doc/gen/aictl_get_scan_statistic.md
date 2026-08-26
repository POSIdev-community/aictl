## aictl get scan statistic

Get scan statistic

### Synopsis

Download or print scan statistics. Scan id comes from argument or stdin. Output path via -o; use -f to overwrite. Use --json for JSON output.

```
aictl get scan statistic <scan-id> [flags]
```

### Examples

```
  aictl get scan statistic <scan-id>
  aictl get scan statistic <scan-id> --json
  aictl get scan statistic <scan-id> -o ./stat.json -f
```

### Options

```
  -f, --force           Overwrite existing output file
  -h, --help            help for statistic
      --json            Output in JSON format
  -o, --output string   Output file path
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

