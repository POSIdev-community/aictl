## aictl scan stop

Stop scan

### Synopsis

Stop a running scan by scan id.

```
aictl scan stop <scan-id> [flags]
```

### Examples

```
  aictl scan stop <scan-id>
```

### Options

```
  -h, --help   help for stop
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

* [aictl scan](aictl_scan.md)	 - Manage scans

