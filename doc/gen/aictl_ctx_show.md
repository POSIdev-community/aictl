## aictl ctx show

Show current aictl context

### Synopsis

Print the local aictl context. Default output is human-readable with a masked token; use --json or --yaml for machine-readable formats with the raw token (mutually exclusive). In JSON/YAML, unset fields are null instead of "<unset>".

```
aictl ctx show [flags]
```

### Examples

```
  aictl ctx show
  aictl ctx show --json
  aictl ctx show --yaml
```

### Options

```
  -h, --help   help for show
      --json   Print context as JSON
      --yaml   Print context as YAML
```

### SEE ALSO

* [aictl ctx](aictl_ctx.md)	 - Manage local aictl context

