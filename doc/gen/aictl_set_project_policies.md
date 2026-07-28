## aictl set project policies

Set project security policies

### Synopsis

Replace project security policies with JSON from an argument, file, or stdin. Project id comes from context or -p.

```
aictl set project policies [json|-] [flags]
```

### Examples

```
  aictl set project policies -f policies.json -p <project-id>
  aictl set project policies -
```

### Options

```
  -f, --file string   Path to policies JSON file, or - for stdin
  -h, --help          help for policies
```

### Options inherited from parent commands

```
  -l, --log-path string     Log file path
  -p, --project-id string   Project id (overrides context)
      --tls-skip            Skip TLS certificate verification
  -t, --token string        AI server access token (overrides context)
  -u, --uri string          AI server URI (overrides context)
  -v, --verbose             Verbose output
```

### SEE ALSO

* [aictl set project](aictl_set_project.md)	 - Set project configuration

