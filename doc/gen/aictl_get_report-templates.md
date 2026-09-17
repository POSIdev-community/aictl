## aictl get report-templates

Get report templates

### Synopsis

List available report templates matching an optional regex filter.

```
aictl get report-templates [<regex>] [flags]
```

### Examples

```
  aictl get report-templates
  aictl get report-templates owasp
  aictl get report-templates -q --localization ru
```

### Options

```
  -h, --help                  help for report-templates
      --localization string   Report localization language: 'en' or 'ru' (default "en")
  -q, --quite                 Print only ids
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

