## aictl get scanning

Get active scans

### Synopsis

List scans currently in progress. Optionally filter by project id; context project id is not used.

```
aictl get scanning [flags]
```

### Examples

```
  aictl get scanning
  aictl get scanning -p <project-id>
```

### Options

```
  -h, --help                help for scanning
  -p, --project-id string   filter by project id (ctx -p is ignored)
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

