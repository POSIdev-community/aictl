## aictl scan project

Start project scan

### Synopsis

Start a scan on an entire project. Project id comes from the argument or context.

```
aictl scan project <project-id> [flags]
```

### Examples

```
  aictl scan project <project-id>
  aictl scan project --scan-label release --full-scan
```

### Options

```
      --full-scan           Run a full scan instead of incremental
  -h, --help                help for project
      --scan-label string   Label for the scan (max 40 chars)
```

### Options inherited from parent commands

```
      --cacert string     Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -V, --debug             Debug output (includes error chains)
  -l, --log-path string   Log file path
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output (operational details)
```

### SEE ALSO

* [aictl scan](aictl_scan.md)	 - Manage scans

