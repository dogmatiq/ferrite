# Environment Variables

| Name               | Usage       | Description             |
| ------------------ | ----------- | ----------------------- |
| [`COLOR_ENABLED`]  | required    | enable colors           |
| [`WIDGET_COLOR`]   | conditional | the color of the widget |
| [`WIDGET_ENABLED`] | required    | enable the widget       |

## `COLOR_ENABLED`

> enable colors

The `COLOR_ENABLED` variable's value **MUST** be either `true` or `false`.

```bash
export COLOR_ENABLED=true
export COLOR_ENABLED=false
```

## `WIDGET_COLOR`

> the color of the widget

The `WIDGET_COLOR` variable **MAY** be left undefined when [`COLOR_ENABLED`] is
`false` or [`WIDGET_ENABLED`] is `false`.

```bash
export WIDGET_COLOR=foo # (non-normative)
```

### See Also

- [`COLOR_ENABLED`] — enable colors
- [`WIDGET_ENABLED`] — enable the widget

## `WIDGET_ENABLED`

> enable the widget

The `WIDGET_ENABLED` variable's value **MUST** be either `true` or `false`.

```bash
export WIDGET_ENABLED=true
export WIDGET_ENABLED=false
```

<!-- references -->

[`color_enabled`]: #COLOR_ENABLED
[`widget_color`]: #WIDGET_COLOR
[`widget_enabled`]: #WIDGET_ENABLED
