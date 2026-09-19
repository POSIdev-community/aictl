## aictl get scan statistic

Get scan statistic

### Synopsis

Print scan statistics (severity counts, files/URLs scanned, duration, policy state).
Scan id comes from argument or stdin. Default output is text; --json prints pretty JSON;
-o writes JSON to a file (use -f to overwrite).
scanDuration is the API ISO-8601 duration (e.g. PT00H34M35.872S); policyState is None, Rejected, or Confirmed.

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
  -V, --debug               Debug output (includes error chains)
  -l, --log-path string     Log file path
  -p, --project-id string   Project id (overrides context)
      --tls-skip            Skip TLS certificate verification
  -t, --token string        AI server access token (overrides context)
  -u, --uri string          AI server URI (overrides context)
  -v, --verbose             Verbose output (operational details)
```

### SEE ALSO

* [aictl get scan](aictl_get_scan.md)	 - Get scan

