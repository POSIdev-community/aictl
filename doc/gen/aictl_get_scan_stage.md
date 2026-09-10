## aictl get scan stage

Get scan stage

### Synopsis

Print the current scan stage. Scan id comes from argument or stdin. Use --fail-on-scan-failed to exit with code 1 when stage is Failed or Aborted.

```
aictl get scan stage <scan-id> [flags]
```

### Examples

```
  aictl get scan stage <scan-id>
  aictl get scan stage <scan-id> --fail-on-scan-failed
```

### Options

```
      --fail-on-scan-failed   Exit with code 1 when scan stage is Failed or Aborted
  -h, --help                  help for stage
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

