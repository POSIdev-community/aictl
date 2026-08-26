## aictl get scan logs

Get scan logs

### Synopsis

Download scan logs to a file. Scan id comes from argument or stdin. Output path -o is required; use -f to overwrite.

```
aictl get scan logs <scan-id> [flags]
```

### Examples

```
  aictl get scan logs <scan-id> -o ./scan.log
  aictl get scan logs <scan-id> -o ./scan.log -f
```

### Options

```
  -f, --force           Overwrite existing output file
  -h, --help            help for logs
  -o, --output string   Output file path
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

