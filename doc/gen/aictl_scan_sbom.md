## aictl scan sbom

Start SBOM scan

### Synopsis

Start a scan on an SBOM project. Project id comes from the argument, -p, or context.

```
aictl scan sbom <project-id> [flags]
```

### Examples

```
  aictl scan sbom <project-id>
  aictl scan sbom -p <project-id> --scan-label release
```

### Options

```
  -h, --help                help for sbom
  -p, --project-id string   Project id (overrides context)
      --scan-label string   Label for the scan (max 40 chars)
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

* [aictl scan](aictl_scan.md)	 - Manage scans

