# Environment Variables

## Specification

### `PRIVATE_KEY`

> path to the private key file

The `PRIVATE_KEY` variable **MAY** be left undefined, in which case the default
value of `/etc/ssh/id_rsa` is used.

```bash
export PRIVATE_KEY=/etc/ssh/id_rsa # (default)
export PRIVATE_KEY=/path/to/file   # (non-normative) an absolute file path
export PRIVATE_KEY=./path/to/file  # (non-normative) a relative file path
```
