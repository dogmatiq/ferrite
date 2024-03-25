# Environment Variables

## `LOG_LEVEL`

> the minimum log level to record

The `LOG_LEVEL` variable **MAY** be left undefined. Otherwise, the value
**MUST** be one of the values shown in the examples below.

```bash
export LOG_LEVEL=debug # show information for developers
export LOG_LEVEL=info  # standard log messages
export LOG_LEVEL=warn  # important, but don't need individual human review
export LOG_LEVEL=error # a healthy application shouldn't produce any errors
export LOG_LEVEL=fatal # the application cannot proceed
```
