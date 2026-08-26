## aictl update project settings

Update project settings

### Synopsis

Patch project scan settings (priority, preferred agents). At least one of --priority, --agents, --preferred-agents-only, or --no-preferred-agents-only is required. Project id comes from context or -p.

```
aictl update project settings [flags]
```

### Examples

```
  aictl update project settings --priority High -p <project-id>
  aictl update project settings --agents agent1,agent2 --preferred-agents-only
```

### Options

```
      --agents string              Comma-separated preferred scan agent ids
  -h, --help                       help for settings
      --no-preferred-agents-only   Allow all scan agents, not only selected ones
      --preferred-agents-only      Use only the selected scan agents
      --priority string            Scan priority: None, Low, Medium, High, or Critical
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

* [aictl update project](aictl_update_project.md)	 - Update project

