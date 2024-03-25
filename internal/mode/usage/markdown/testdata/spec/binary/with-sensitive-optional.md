# Environment Variables

## `SECRET_KEY`

> a very secret machine-readable key

The `SECRET_KEY` variable **MAY** be left undefined. Otherwise, the value
**MUST** be a binary value expressed using the `base64` encoding scheme.

⚠️ This variable is **sensitive**; its value may contain private information.
