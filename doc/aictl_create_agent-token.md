## aictl create agent-token

Create access token for scan agent

### Synopsis

Create access token for AIE scan agent.

Requires admin credentials (login/password) to authenticate and create
a token with ScanAgent scope. The generated token can be used to
configure scan agents.

Example:
  aictl create agent-token my-agent --login admin --password secret -u https://aie-server:443

```
aictl create agent-token <agent-name> [flags]
```

### Options

```
  -h, --help              help for agent-token
      --login string      admin user login (required)
      --password string   admin user password (required)
```

### Options inherited from parent commands

```
  -l, --log-path string   log file path
      --safe              if resource exists, return its id without error
      --tls-skip          Skip certificate verification
  -t, --token string      AI server access token
  -u, --uri string        AI server uri
  -v, --verbose           verbose output
```

### SEE ALSO

* [aictl create](aictl_create.md)	 - Create resource

