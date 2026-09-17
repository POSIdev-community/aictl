## aictl update

Update resources

### Synopsis

Update project sources, SBOM files, SCA feeds, and settings on the server.

### Options

```
      --cacert string     Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -V, --debug             Debug output (includes error chains)
  -h, --help              help for update
  -l, --log-path string   Log file path
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output (operational details)
```

### SEE ALSO

* [aictl](aictl.md)	 - Application Inspector ConTroL tool
* [aictl update project](aictl_update_project.md)	 - Update project
* [aictl update sbom](aictl_update_sbom.md)	 - Update SBOM
* [aictl update sca-feeds](aictl_update_sca-feeds.md)	 - Upload SCA feeds package
* [aictl update sources](aictl_update_sources.md)	 - Update sources

