## aictl get project exclusions

Get project file/folder exclusions

### Synopsis

Print project file and folder exclusions (gitignore-style). Project id comes from context or parent -p.

```
aictl get project exclusions [flags]
```

### Examples

```
  aictl get project exclusions -p <project-id>
```

### Options

```
  -h, --help   help for exclusions
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

* [aictl get project](aictl_get_project.md)	 - Get project

