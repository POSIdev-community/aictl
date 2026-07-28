## aictl update project settings

Update project settings

```
aictl update project settings [flags]
```

### Options

```
      --agents string              comma-separated scan agent ids
  -h, --help                       help for settings
      --no-preferred-agents-only   allow all agents, not only selected ones
      --preferred-agents-only      use only selected scan agents
      --priority string            scan priority: None, Low, Medium, High, Critical
```

### Options inherited from parent commands

```
  -l, --log-path string     log file path
  -p, --project-id string   project id
      --tls-skip            Skip certificate verification
  -t, --token string        AI server access token
  -u, --uri string          AI server uri
  -v, --verbose             verbose output
```

### SEE ALSO

* [aictl update project](aictl_update_project.md)	 - Update project

