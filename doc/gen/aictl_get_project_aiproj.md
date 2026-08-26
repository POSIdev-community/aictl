## aictl get project aiproj

Get project aiproj

### Synopsis

Download the project .aiproj configuration file. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.

```
aictl get project aiproj [flags]
```

### Examples

```
  aictl get project aiproj -p <project-id> -o ./project.aiproj
  aictl get project aiproj -p <project-id> -o ./project.aiproj -f
```

### Options

```
  -f, --force           Overwrite existing output file
  -h, --help            help for aiproj
  -o, --output string   Output file path
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

