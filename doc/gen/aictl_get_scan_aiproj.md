## aictl get scan aiproj

Get scan aiproj

### Synopsis

Download the .aiproj configuration used for the scan. Scan id comes from argument or stdin. Output path via -o; use -f to overwrite.

```
aictl get scan aiproj <scan-id> [flags]
```

### Examples

```
  aictl get scan aiproj <scan-id> -o ./scan.aiproj
  aictl get scan aiproj <scan-id> -o ./scan.aiproj -f
```

### Options

```
  -f, --force           Overwrite existing output file
  -h, --help            help for aiproj
  -o, --output string   Output file path
```

### Options inherited from parent commands

```
  -l, --log-path string     Log file path
  -p, --project-id string   Project id (overrides context)
      --tls-skip            Skip TLS certificate verification
  -t, --token string        AI server access token (overrides context)
  -u, --uri string          AI server URI (overrides context)
  -v, --verbose             Verbose output
```

### SEE ALSO

* [aictl get scan](aictl_get_scan.md)	 - Get scan

