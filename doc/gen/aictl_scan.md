## aictl scan

Manage scans

### Synopsis

Start, stop, await, and check security scans on the server.

### Options

```
      --cacert string     Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -h, --help              help for scan
  -l, --log-path string   Log file path
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output
```

### SEE ALSO

* [aictl](aictl.md)	 - Application Inspector ConTroL tool
* [aictl scan await](aictl_scan_await.md)	 - Await scan completion
* [aictl scan check-policies](aictl_scan_check-policies.md)	 - Check scan policy state
* [aictl scan start](aictl_scan_start.md)	 - Start scan
* [aictl scan stop](aictl_scan_stop.md)	 - Stop scan

