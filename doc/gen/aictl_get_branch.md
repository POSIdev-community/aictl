## aictl get branch

Get branch

### Synopsis

Retrieve branch details by id. Branch id may be passed as an argument or via stdin.

```
aictl get branch <branch-id> [flags]
```

### Examples

```
  aictl get branch <branch-id>
```

### Options

```
  -h, --help   help for branch
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

