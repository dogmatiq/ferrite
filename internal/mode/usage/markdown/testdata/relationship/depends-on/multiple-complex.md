# Environment Variables

| Name              | Usage       | Description             |
| ----------------- | ----------- | ----------------------- |
| [`COLOR_ENABLED`] | required    | enable colors           |
| [`WIDGET_COLOR`]  | conditional | the color of the widget |
| [`WIDGET_MODE`]   | required    | set the widget mode     |

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
`false` or [`WIDGET_MODE`] is not `color`.

```bash
export WIDGET_COLOR=foo # (non-normative)
```

### See Also

- [`COLOR_ENABLED`] — enable colors
- [`WIDGET_MODE`] — set the widget mode

## `WIDGET_MODE`

> set the widget mode

The `WIDGET_MODE` variable's value **MUST** be either `grayscale` or `color`.

```bash
export WIDGET_MODE=grayscale
export WIDGET_MODE=color
```

<!-- references -->

[`color_enabled`]: #COLOR_ENABLED
[`widget_color`]: #WIDGET_COLOR
[`widget_mode`]: #WIDGET_MODE
