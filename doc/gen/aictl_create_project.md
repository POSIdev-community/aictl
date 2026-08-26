## aictl create project

Create a project

### Synopsis

Create a new AI project by name. The name may be passed as an argument or via stdin. With --safe, an existing project id is returned instead of an error.

```
aictl create project <project-name> [flags]
```

### Examples

```
  aictl create project my-app
  aictl create project my-app --safe
  echo my-app | aictl create project
```

### Options

```
  -h, --help   help for project
```

### Options inherited from parent commands

```
      --cacert string     Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -l, --log-path string   Log file path
      --safe              If the resource already exists, return its id without error
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output
```

### SEE ALSO

* [aictl create](aictl_create.md)	 - Create resources

