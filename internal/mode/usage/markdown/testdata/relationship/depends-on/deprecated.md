# Environment Variables

| Name                 | Optionality          | Description             |
| -------------------- | -------------------- | ----------------------- |
| ~~[`WIDGET_COLOR`]~~ | optional, deprecated | the color of the widget |
| [`WIDGET_ENABLED`]   | required             | enable the widget       |

## Specification

### `WIDGET_COLOR`

> the color of the widget

⚠️ The `WIDGET_COLOR` variable is **deprecated**; its use is **NOT RECOMMENDED**
as it may be removed in a future version. The value is not used when
[`WIDGET_ENABLED`] is `false`.

```bash
export WIDGET_COLOR=foo # (non-normative)
```

#### See Also

- [`WIDGET_ENABLED`] — enable the widget

### `WIDGET_ENABLED`

> enable the widget

The `WIDGET_ENABLED` variable's value **MUST** be either `true` or `false`.

```bash
export WIDGET_ENABLED=true
export WIDGET_ENABLED=false
```

<!-- references -->

[`widget_color`]: #WIDGET_COLOR
[`widget_enabled`]: #WIDGET_ENABLED
