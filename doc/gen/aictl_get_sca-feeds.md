## aictl get sca-feeds

List or download SCA feeds packages

### Synopsis

List SCA feeds packages or download one by version (AIE ≥ 6.3).

Without arguments — print table. With <version> — download archive.
Optional -o: new file path, or existing directory (writes default file name inside).

```
aictl get sca-feeds [version] [flags]
```

### Examples

```
  aictl get sca-feeds
  aictl get sca-feeds --status current --status active
  aictl get sca-feeds 1.2.3
  aictl get sca-feeds 1.2.3 -o ./feeds.zip
  aictl get sca-feeds 1.2.3 -o ./outdir/
```

### Options

```
  -h, --help                 help for sca-feeds
  -o, --output string        Output file path or existing directory
      --status stringArray   Filter by status (repeatable): active, current, archived, rolled_back
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

