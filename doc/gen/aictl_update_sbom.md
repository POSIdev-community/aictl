## aictl update sbom

Update SBOM

### Synopsis

Upload or replace the SBOM file for an SBOM project. Project id comes from context or -p.

```
aictl update sbom <path> [flags]
```

### Examples

```
  aictl update sbom ./sbom.json -p <project-id>
  aictl update sbom ./bom.xml
```

### Options

```
  -h, --help                help for sbom
  -p, --project-id string   Project id (overrides context)
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

* [aictl update](aictl_update.md)	 - Update resources

