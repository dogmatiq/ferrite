# Environment Variables

| Name               | Usage    | Description             |
| ------------------ | -------- | ----------------------- |
| [`WIDGET_COLOR`]   | optional | the color of the widget |
| [`WIDGET_ENABLED`] | required | enable the widget       |

## `WIDGET_COLOR`

> the color of the widget

The `WIDGET_COLOR` variable **MAY** be left undefined. It is ignored when
[`WIDGET_ENABLED`] is `false`.

```bash
export WIDGET_COLOR=foo # (non-normative)
```

### See Also

- [`WIDGET_ENABLED`] — enable the widget

## `WIDGET_ENABLED`

> enable the widget

The `WIDGET_ENABLED` variable's value **MUST** be either `true` or `false`.

```bash
export WIDGET_ENABLED=true
export WIDGET_ENABLED=false
```

<!-- references -->

[`widget_color`]: #WIDGET_COLOR
[`widget_enabled`]: #WIDGET_ENABLED
