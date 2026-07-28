## aictl scan await

Await scan completion

### Synopsis

Wait until a scan reaches a terminal stage. Project id comes from context or -p.

```
aictl scan await <scan-id> [flags]
```

### Examples

```
  aictl scan await <scan-id>
  aictl scan await <scan-id> -p <project-id> --fail-on-scan-failed
```

### Options

```
      --fail-on-scan-failed   Exit 1 if scan stage is Failed or Aborted
  -h, --help                  help for await
  -p, --project-id string     Project id (overrides context)
```

### Options inherited from parent commands

```
  -l, --log-path string   Log file path
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output
```

### SEE ALSO

* [aictl scan](aictl_scan.md)	 - Manage scans

