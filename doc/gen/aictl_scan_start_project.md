## aictl scan start project

Start project scan

### Synopsis

Start a scan on an entire project. Project id comes from the argument or context.

```
aictl scan start project <project-id> [flags]
```

### Examples

```
  aictl scan start project <project-id>
  aictl scan start project --scan-label release --full-scan
```

### Options

```
  -h, --help   help for project
```

### Options inherited from parent commands

```
      --full-scan           Run a full scan instead of incremental
  -l, --log-path string     Log file path
      --scan-label string   Label for the scan (max 40 chars)
      --tls-skip            Skip TLS certificate verification
  -t, --token string        AI server access token (overrides context)
  -u, --uri string          AI server URI (overrides context)
  -v, --verbose             Verbose output
```

### SEE ALSO

* [aictl scan start](aictl_scan_start.md)	 - Start scan

