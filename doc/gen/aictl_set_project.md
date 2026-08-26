## aictl set project

Set project configuration

### Synopsis

Set project-level configuration on the server. Project id comes from context or -p.

### Options

```
  -h, --help                help for project
  -p, --project-id string   Project id (overrides context)
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

* [aictl set](aictl_set.md)	 - Set resource configuration
* [aictl set project exclusions](aictl_set_project_exclusions.md)	 - Set project exclusions
* [aictl set project policies](aictl_set_project_policies.md)	 - Set project security policies
* [aictl set project settings](aictl_set_project_settings.md)	 - Set project settings

