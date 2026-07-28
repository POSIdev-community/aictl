## aictl get branches

Get AI branches

### Synopsis

List branches for the current project matching an optional regex filter. Project id comes from context or -p.

```
aictl get branches [<regex>] [flags]
```

### Examples

```
  aictl get branches -p <project-id>
  aictl get branches main -p <project-id>
  aictl get branches -p <project-id> -q
```

### Options

```
  -h, --help                help for branches
  -p, --project-id string   Project id (overrides context)
  -q, --quite               Print only ids
```

### Options inherited from parent commands

```
  -l, --log-path string   Log file path
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output
```

### SEE ALSO

* [aictl get](aictl_get.md)	 - Get resources

