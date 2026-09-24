## aictl get scans sbom-project

Get scans for an SBOM project

### Synopsis

List scans for an SBOM project matching an optional regex filter. Resolves the virtual branch from the project id (context or -p).

```
aictl get scans sbom-project [<regex>] [flags]
```

### Examples

```
  aictl get scans sbom-project -p <project-id>
  aictl get scans sbom-project 2024 -p <project-id>
  aictl get scans sbom-project --latest -p <project-id>
```

### Options

```
  -h, --help                help for sbom-project
      --latest              Return only the latest scan result
  -p, --project-id string   Project id (overrides context)
  -q, --quiet               Print only ids
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

* [aictl get scans](aictl_get_scans.md)	 - Get AI scans

