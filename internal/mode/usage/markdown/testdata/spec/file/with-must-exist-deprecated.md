# Environment Variables

## `PRIVATE_KEY`

> path to the private key file

⚠️ The `PRIVATE_KEY` variable is **deprecated**; its use is **NOT RECOMMENDED**
as it may be removed in a future version. If defined, the value **MUST** refer
to a file that already exists.

```bash
export PRIVATE_KEY=/path/to/file  # (non-normative) an absolute file path
export PRIVATE_KEY=./path/to/file # (non-normative) a relative file path
```
