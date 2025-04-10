# Environment Variables

| Name             | Usage       | Description             |
| ---------------- | ----------- | ----------------------- |
| [`WIDGET_COLOR`] | conditional | the color of the widget |
| [`WIDGET_TYPE`]  | required    | the type of widget      |

## `WIDGET_COLOR`

> the color of the widget

The `WIDGET_COLOR` variable **MAY** be left undefined when [`WIDGET_TYPE`] is
undefined.

```bash
export WIDGET_COLOR=foo # (non-normative)
```

### See Also

- [`WIDGET_TYPE`] — the type of widget

## `WIDGET_TYPE`

> the type of widget

The `WIDGET_TYPE` variable **MUST NOT** be left undefined.

```bash
export WIDGET_TYPE=foo # (non-normative)
```

<!-- references -->

[`widget_color`]: #WIDGET_COLOR
[`widget_type`]: #WIDGET_TYPE
