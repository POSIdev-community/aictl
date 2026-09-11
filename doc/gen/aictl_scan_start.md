## aictl scan start

Start scan (deprecated)

### Synopsis

Deprecated: use 'aictl scan branch' or 'aictl scan project'. Start a full or incremental scan on a project or branch.

```
aictl scan start [flags]
```

### Options

```
      --full-scan           Run a full scan instead of incremental
  -h, --help                help for start
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
* [aictl scan start branch](aictl_scan_start_branch.md)	 - Start branch scan (deprecated)
* [aictl scan start project](aictl_scan_start_project.md)	 - Start project scan (deprecated)

