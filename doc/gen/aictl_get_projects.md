## aictl get projects

Get AI projects

### Synopsis

List projects matching an optional regex filter. Filter may be passed as an argument or via stdin.

```
aictl get projects <regex> [flags]
```

### Examples

```
  aictl get projects
  aictl get projects my-.*
  aictl get projects -q
```

### Options

```
  -h, --help    help for projects
  -q, --quite   Print only ids
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

* [aictl get](aictl_get.md)	 - Get resources

