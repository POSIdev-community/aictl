## aictl set project settings

Set project settings

### Synopsis

Replace project settings (aiproj.json) from a JSON argument or file. Project id comes from context or -p.

```
aictl set project settings [flags]
```

### Examples

```
  aictl set project settings -f aiproj.json -p <project-id>
  aictl set project settings '{"Version":"1.0"}' -p <project-id>
```

### Options

```
  -f, --file string   Path to aiproj.json file
  -h, --help          help for settings
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

