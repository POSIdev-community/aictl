## aictl update project languages

Update project languages

### Synopsis

Detect and update project languages from uploaded sources. Project id comes from context or -p.

```
aictl update project languages [flags]
```

### Examples

```
  aictl update project languages -p <project-id>
```

### Options

```
  -h, --help   help for languages
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

* [aictl update project](aictl_update_project.md)	 - Update project

