## aictl update sca-feeds

Upload SCA feeds package

### Synopsis

Upload a SCA feeds zip archive to the server (AIE ≥ 6.3). Requires --version.

```
aictl update sca-feeds <path> [flags]
```

### Examples

```
  aictl update sca-feeds ./AI.SCA.Feeds.47.zip --version 47
  aictl update sca-feeds ./feeds.zip --version 1.2.3
```

### Options

```
  -h, --help             help for sca-feeds
      --version string   Package version (required)
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

* [aictl update](aictl_update.md)	 - Update resources

