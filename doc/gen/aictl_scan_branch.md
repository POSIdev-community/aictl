## aictl scan branch

Start branch scan

### Synopsis

Start a scan on a branch. Branch id comes from the argument or context; project id from context or -p.

```
aictl scan branch <branch-id> [flags]
```

### Examples

```
  aictl scan branch <branch-id> -p <project-id>
  aictl scan branch --scan-label nightly --full-scan
```

### Options

```
      --full-scan           Run a full scan instead of incremental
  -h, --help                help for branch
  -p, --project-id string   Project id (overrides context)
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

