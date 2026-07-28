## aictl

Application Inspector ConTroL tool

### Synopsis

CLI for managing PT Application Inspector: projects, branches, scans, reports, and local context.

Use subcommands to talk to an AI server. Connection settings come from context (~/.config/aictl/context.yaml) or -u/-t/--tls-skip flags.

```
aictl [flags]
```

### Examples

```
  aictl ctx set -u https://ai.example -t <token>
  aictl get projects
  aictl --version
```

### Options

```
  -h, --help      help for aictl
  -v, --version   Show aictl version
```

### SEE ALSO

* [aictl create](aictl_create.md)	 - Create resources
* [aictl ctx](aictl_ctx.md)	 - Manage local aictl context
* [aictl delete](aictl_delete.md)	 - Delete resources
* [aictl get](aictl_get.md)	 - Get resources
* [aictl scan](aictl_scan.md)	 - Manage scans
* [aictl set](aictl_set.md)	 - Set resource configuration
* [aictl update](aictl_update.md)	 - Update resources

