## aictl get scan sbom

Get scan SBOM

### Synopsis

Download the scan SBOM (Software Bill of Materials). Scan id comes from argument or stdin. Output path via -o; use -f to overwrite.

```
aictl get scan sbom <scan-id> [flags]
```

### Examples

```
  aictl get scan sbom <scan-id> -o ./sbom.json
  aictl get scan sbom <scan-id> -o ./sbom.json -f
```

### Options

```
  -f, --force           Overwrite existing output file
  -h, --help            help for sbom
  -o, --output string   Output file path
```

### Options inherited from parent commands

```
      --cacert string       Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -V, --debug               Debug output (includes error chains)
  -l, --log-path string     Log file path
  -p, --project-id string   Project id (overrides context)
      --tls-skip            Skip TLS certificate verification
  -t, --token string        AI server access token (overrides context)
  -u, --uri string          AI server URI (overrides context)
  -v, --verbose             Verbose output (operational details)
```

### SEE ALSO

* [aictl get scan](aictl_get_scan.md)	 - Get scan

