## aictl create sbom-project

Create an SBOM project

### Synopsis

Create a new SBOM project and upload an SBOM file. The name may be passed as an argument or via stdin. With --safe, an existing SBOM project is updated instead of an error.

```
aictl create sbom-project <project-name> [flags]
```

### Examples

```
  aictl create sbom-project my-sbom --file ./sbom.json
  aictl create sbom-project my-sbom --file ./sbom.json --safe
  echo my-sbom | aictl create sbom-project - --file ./sbom.json
```

### Options

```
      --file string   Path to SBOM file to upload
  -h, --help          help for sbom-project
```

### Options inherited from parent commands

```
      --cacert string     Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -V, --debug             Debug output (includes error chains)
  -l, --log-path string   Log file path
      --safe              If the resource already exists, return its id without error
      --tls-skip          Skip TLS certificate verification
  -t, --token string      AI server access token (overrides context)
  -u, --uri string        AI server URI (overrides context)
  -v, --verbose           Verbose output (operational details)
```

### SEE ALSO

* [aictl create](aictl_create.md)	 - Create resources

