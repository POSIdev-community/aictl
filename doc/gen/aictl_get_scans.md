## aictl get scans

Get AI scans

### Synopsis

List scans for the current branch matching an optional regex filter. Branch id comes from context or -b.

```
aictl get scans [<regex>] [flags]
```

### Examples

```
  aictl get scans -b <branch-id>
  aictl get scans 2024 -b <branch-id>
  aictl get scans --latest -b <branch-id>
```

### Options

```
  -b, --branch-id string   Branch id (overrides context)
  -h, --help               help for scans
      --latest             Return only the latest scan result
  -q, --quite              Print only ids
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

* [aictl get](aictl_get.md)	 - Get resources

