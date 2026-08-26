## aictl ctx unset

Unset context parameters

### Synopsis

Clear selected fields from the local aictl context. At least one flag is required.

```
aictl ctx unset [flags]
```

### Examples

```
  aictl ctx unset -p -b
  aictl ctx unset -u -t
  aictl ctx unset --cacert
```

### Options

```
  -b, --branch-id    Unset branch id
      --cacert       Unset CA certificate path
  -h, --help         help for unset
  -p, --project-id   Unset project id
      --tls-skip     Unset TLS skip setting
  -t, --token        Unset access token
  -u, --uri          Unset URI
```

### SEE ALSO

* [aictl ctx](aictl_ctx.md)	 - Manage local aictl context

