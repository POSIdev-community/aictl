## aictl ctx set

Set current aictl configuration

### Synopsis

Update one or more fields in the local aictl context. At least one flag is required. Do not pass both --tls-skip and --no-tls-skip. Do not pass --cacert together with --tls-skip. If the other option is already stored in context, unset it first (aictl ctx unset --tls-skip or aictl ctx unset --cacert).

```
aictl ctx set [flags]
```

### Examples

```
  aictl ctx set -u https://ai.example -t <token>
  aictl ctx set -p <project-id> -b <branch-id>
  aictl ctx set --cacert /path/to/ca.pem
  aictl ctx set --tls-skip
```

### Options

```
  -b, --branch-id string    Default branch id
      --cacert string       Path to PEM file with CA certificate(s) to trust (appended to system roots)
  -h, --help                help for set
      --no-tls-skip         Require TLS certificate verification
  -p, --project-id string   Default project id
      --tls-skip            Skip TLS certificate verification
  -t, --token string        AI server access token
  -u, --uri string          AI server URI
```

### SEE ALSO

* [aictl ctx](aictl_ctx.md)	 - Manage local aictl context

