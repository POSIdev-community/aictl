## aictl get agents

Get agents

### Synopsis

List scan agents registered on the server.

```
aictl get agents [flags]
```

### Examples

```
  aictl get agents
  aictl get agents -q
```

### Options

```
  -h, --help    help for agents
  -q, --quiet   Print only ids
```

### Options inherited from parent commands

```
      --cacert string     Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -V, --debug             Debug output (includes error chains)
  -l, --log-path string   Log file path
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output (operational details)
```

### SEE ALSO

* [aictl get](aictl_get.md)	 - Get resources

