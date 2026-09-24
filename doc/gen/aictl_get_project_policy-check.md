## aictl get project policy-check

Get project security policy check flag

### Synopsis

Print checkSecurityPoliciesAccordance as true or false. Project id comes from context or parent -p.

```
aictl get project policy-check [flags]
```

### Examples

```
  aictl get project policy-check -p <project-id>
```

### Options

```
  -h, --help   help for policy-check
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

