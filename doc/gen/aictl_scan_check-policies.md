## aictl scan check-policies

Check scan policy state

### Synopsis

Check the policy state for a scan. Project id comes from context or -p.

```
aictl scan check-policies <scan-id> [flags]
```

### Examples

```
  aictl scan check-policies <scan-id>
  aictl scan check-policies <scan-id> --fail-on-policies-rejected
```

### Options

```
      --fail-on-policies-rejected   Exit 1 if PolicyState is Rejected
  -h, --help                        help for check-policies
  -p, --project-id string           Project id (overrides context)
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

