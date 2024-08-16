# Environment Variables

## `CACHE_DIR`

> path to the cache directory

The `CACHE_DIR` variable **MAY** be left undefined, in which case the default
value of `/var/run/cache` is used. Otherwise, the value **MUST** refer to a
directory that already exists.

```bash
export CACHE_DIR=/var/run/cache # (default)
export CACHE_DIR=/path/to/dir   # (non-normative) an absolute directory path
export CACHE_DIR=./path/to/dir  # (non-normative) a relative directory path
```
