## aictl update sca-feeds rollback

Roll back current SCA feeds package

### Synopsis

Roll back the current SCA feeds package to the previous version selected by the server (AIE ≥ 6.3). Prints the new version on stdout.

```
aictl update sca-feeds rollback [flags]
```

### Examples

```
  aictl update sca-feeds rollback
  aictl update sca-feeds rollback -y
```

### Options

```
  -h, --help   help for rollback
  -y, --yes    Skip confirmation prompt
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

* [aictl update sca-feeds](aictl_update_sca-feeds.md)	 - Upload SCA feeds package

