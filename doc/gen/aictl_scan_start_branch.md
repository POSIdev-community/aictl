## aictl scan start branch

Start branch scan

### Synopsis

Start a scan on a branch. Branch id comes from the argument or context; project id from context or -p.

```
aictl scan start branch <branch-id> [flags]
```

### Examples

```
  aictl scan start branch <branch-id> -p <project-id>
  aictl scan start branch --scan-label nightly --full-scan
```

### Options

```
  -h, --help                help for branch
  -p, --project-id string   Project id (overrides context)
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

