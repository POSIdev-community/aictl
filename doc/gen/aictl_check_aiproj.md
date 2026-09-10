## aictl check aiproj

Check aiproj schema

### Synopsis

Validate an aiproj/JSON document against JSON Schema offline. Input via -f, positional JSON, or stdin (same as set project settings).

```
aictl check aiproj [flags]
```

### Examples

```
  aictl check aiproj -f aiproj.json
  aictl check aiproj -f aiproj.json --schema-version 1.11
  aictl check aiproj -f aiproj.json --json
  aictl check aiproj '{"Version":"1.11","ProjectName":"demo","ProgrammingLanguages":["Go"],"ScanModules":["StaticCodeAnalysis"]}'
```

### Options

```
  -f, --file string             Path to aiproj.json file
  -h, --help                    help for aiproj
      --json                    Print result as pretty JSON on stdout
      --schema-version string   Schema version to validate against (1.8, 1.9, 1.10, 1.11); default: auto-detect
```

### Options inherited from parent commands

```
  -V, --debug             Debug output (includes error chains)
  -l, --log-path string   Log file path
  -v, --verbose           Verbose output (operational details)
```

### SEE ALSO

* [aictl check](aictl_check.md)	 - Validate local resources

