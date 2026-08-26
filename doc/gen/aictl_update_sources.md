## aictl update sources

Update sources

### Synopsis

Upload or update sources for a project/branch with optional gitignore-style exclusions. Project and branch ids come from context or -p/-b.

```
aictl update sources <path> [flags]
```

### Examples

```
  aictl update sources ./src -p <project-id> -b <branch-id>
  aictl update sources ./src -e '*.tmp' --exclude-from .aictlignore
```

### Options

```
  -b, --branch-id string           Branch id (overrides context)
  -e, --exclude stringArray        Exclude path (gitignore pattern); repeatable
      --exclude-from stringArray   File with gitignore-style exclude patterns
  -h, --help                       help for sources
  -p, --project-id string          Project id (overrides context)
      --temp-dir string            Directory for temporary zip when packing sources
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

* [aictl update](aictl_update.md)	 - Update resources

