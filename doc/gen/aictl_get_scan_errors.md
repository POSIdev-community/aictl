## aictl get scan errors

Get scan errors

### Synopsis

Print scan errors. Scan id comes from argument or stdin.

```
aictl get scan errors <scan-id> [flags]
```

### Examples

```
  aictl get scan errors <scan-id>
```

### Options

```
  -h, --help   help for errors
```

### Options inherited from parent commands

```
      --cacert string       Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -l, --log-path string     Log file path
  -p, --project-id string   Project id (overrides context)
      --tls-skip            Skip TLS certificate verification
  -t, --token string        AI server access token (overrides context)
  -u, --uri string          AI server URI (overrides context)
  -v, --verbose             Verbose output
```

### SEE ALSO

* [aictl get scan](aictl_get_scan.md)	 - Get scan

