## aictl get project policies

Get project security policies

### Synopsis

Print project security policies. Project id comes from context or parent -p.

```
aictl get project policies [flags]
```

### Examples

```
  aictl get project policies -p <project-id>
```

### Options

```
  -h, --help   help for policies
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

* [aictl get project](aictl_get_project.md)	 - Get project

