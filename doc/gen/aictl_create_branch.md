## aictl create branch

Create a branch

### Synopsis

Create a branch under a project. Optionally pack and upload sources from --scan-target with gitignore-style exclusions. Project id comes from context or -p.

```
aictl create branch <branch-name> [flags]
```

### Examples

```
  aictl create branch main -p <project-id>
  aictl create branch main -p <project-id> -s ./src -e '*.tmp' --exclude-from .aictlignore
  aictl create branch main --safe
```

### Options

```
  -e, --exclude stringArray        Exclude path (gitignore pattern); repeatable
      --exclude-from stringArray   File with gitignore-style exclude patterns
  -h, --help                       help for branch
  -p, --project-id string          Project id (overrides context)
  -s, --scan-target string         Path to sources to pack and upload
      --temp-dir string            Directory for temporary zip when packing sources
```

### Options inherited from parent commands

```
      --cacert string     Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -V, --debug             Debug output (includes error chains)
  -l, --log-path string   Log file path
      --safe              If the resource already exists, return its id without error
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output (operational details)
```

### SEE ALSO

* [aictl create](aictl_create.md)	 - Create resources

