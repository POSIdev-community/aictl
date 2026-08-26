## aictl get scan

Get scan

### Synopsis

Retrieve scan-level data. Scan id may be passed as an argument or via stdin. Project id comes from context or -p.

```
aictl get scan <scan-id> [flags]
```

### Options

```
  -h, --help                help for scan
  -p, --project-id string   Project id (overrides context)
```

### Options inherited from parent commands

```
      --cacert string     Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -l, --log-path string   Log file path
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output
```

### SEE ALSO

* [aictl get](aictl_get.md)	 - Get resources
* [aictl get scan aiproj](aictl_get_scan_aiproj.md)	 - Get scan aiproj
* [aictl get scan errors](aictl_get_scan_errors.md)	 - Get scan errors
* [aictl get scan logs](aictl_get_scan_logs.md)	 - Get scan logs
* [aictl get scan report](aictl_get_scan_report.md)	 - Get scan report
* [aictl get scan sbom](aictl_get_scan_sbom.md)	 - Get scan SBOM
* [aictl get scan stage](aictl_get_scan_stage.md)	 - Get scan stage
* [aictl get scan statistic](aictl_get_scan_statistic.md)	 - Get scan statistic

