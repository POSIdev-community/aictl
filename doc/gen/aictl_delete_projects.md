## aictl delete projects

Delete AI projects

### Synopsis

Delete one or more projects by id. Ids may be passed as arguments or via stdin.

```
aictl delete projects <project-id>... [flags]
```

### Examples

```
  aictl delete projects <project-id>
  aictl delete projects <id1> <id2>
```

### Options

```
  -h, --help   help for projects
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

* [aictl delete](aictl_delete.md)	 - Delete resources

