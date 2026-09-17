## aictl create

Create resources

### Synopsis

Create AI projects (source or SBOM) and branches on the server.

### Options

```
      --cacert string     Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -V, --debug             Debug output (includes error chains)
  -h, --help              help for create
  -l, --log-path string   Log file path
      --safe              If the resource already exists, return its id without error
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output (operational details)
```

### SEE ALSO

* [aictl](aictl.md)	 - Application Inspector ConTroL tool
* [aictl create branch](aictl_create_branch.md)	 - Create a branch
* [aictl create project](aictl_create_project.md)	 - Create a project
* [aictl create sbom-project](aictl_create_sbom-project.md)	 - Create an SBOM project

