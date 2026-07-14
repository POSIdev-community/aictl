## aictl
PT Application Inspector CLI control
___
### Documentation

- [Command reference](doc/aictl.md)
- [Migration: aisa → aictl](doc/migration/aisa-to-aictl.md)
- [Migration: ptai-cli-plugin → aictl](doc/migration/ptai-cli-plugin-to-aictl.md)
- [Gap analysis](doc/migration/gap-analysis.md)

### E2E tests

End-to-end tests run the base pipeline from [`examples/base-pipeline.sh`](examples/base-pipeline.sh) against external AIE stands (5.4 and 6.0).

**Setup:**

```bash
make e2e-config
# Edit tests/e2e/stands.local.yaml — set url and token for each stand
make test-e2e
```

**Run one stand:**

```bash
go test -tags=e2e -v -timeout=45m -run 'AIE_5\.4' ./tests/e2e/...
```

The test builds `bin/aictl` before the pipeline (same as `make build-e2e`). Override with `AICTL_BIN=/path/to/aictl` to skip the build. Config path: `AICTL_E2E_CONFIG=/path/to/stands.yaml`.

Without `stands.local.yaml`, e2e tests are skipped; `go test ./...` does not run them (build tag `e2e`).
