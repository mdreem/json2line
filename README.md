# json2line

`json2line` is a tool that takes JSON strings in lines and formats its values to a string.
The format can be configured with Go templates.

This can be done as follows:

```bash
echo { "key": "value" } | json2-line -f templates.toml -f "my_template"
```

with a config file

```toml
my_template = "the value is:'{{ .key }}'"
```

will output

```
the value is:'value'
```
