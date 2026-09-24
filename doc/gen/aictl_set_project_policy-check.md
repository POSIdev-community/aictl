## aictl set project policy-check

Enable or disable project security policy check

### Synopsis

Set checkSecurityPoliciesAccordance for the project (UI "use security policies" checkbox).
Existing policy rules are preserved. Project id comes from context or -p.

```
aictl set project policy-check <true|false> [flags]
```

### Examples

```
  aictl set project policy-check true
  aictl set project policy-check false -p <project-id>
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

* [aictl set project](aictl_set_project.md)	 - Set project configuration

