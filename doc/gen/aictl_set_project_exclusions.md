## aictl set project exclusions

Set project exclusions

### Synopsis

Replace project file/folder exclusions with gitignore-like text from an argument, file, or stdin. Existing exclusions are fully replaced.

```
aictl set project exclusions [text|-] [flags]
```

### Examples

```
  aictl set project exclusions -f .aictlignore -p <project-id>
  aictl set project exclusions -
  aictl set project exclusions '*.tmp' -p <project-id>
```

### Options

```
  -f, --file string   Path to exclusions file, or - for stdin
  -h, --help          help for exclusions
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

