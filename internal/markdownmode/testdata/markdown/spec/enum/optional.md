# Environment Variables

## Specification

### `LOG_LEVEL`

> the minimum log level to record

This variable **MAY** be set to one of the values below or left undefined.

```bash
export LOG_LEVEL=debug # show information for developers
export LOG_LEVEL=info  # standard log messages
export LOG_LEVEL=warn  # important, but don't need individual human review
export LOG_LEVEL=error # a healthy application shouldn't produce any errors
export LOG_LEVEL=fatal # the application cannot proceed
```
