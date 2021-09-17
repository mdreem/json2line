# json2line

`json2line` is a tool that takes JSON strings in lines and formats its values to a string.
The format can be configured with Go templates.

This can be done as follows:

```bash
echo { "key": "value" } | json2-line -f templates.toml -f "my_template"
```

with a config file

```toml
[templates]
my_template = "the value is:'{{ .key }}'"
```

will output

```
the value is:'value'
```

## Special characters

If you need to replace keys containing special characters, like e.g. `@` a replacement can be defined
as follows:

```toml
[templates]
my_template = "the value is:'{{ .at_key }}'"
[replacements]
"@"="at_"
```

Now all occurrences of `@` will be replaced with `at_`. Now, if `@key` appears in the JSON, it will
be renamed to `at_key` internally so that it can be accessed via `{{ .at_key }}` as shown in the
`templates`-section.

## Configuration

The configuration can be adapted via the command line interface

This can be achieved by the subcommand `json2line configure` which has two subcommands `formatter` and `replacement`
which can be used to change the available formatters and replacements, respectively.

### Changing the formatter configuration

```bash
# add a new formatter with the name 'new_formatter'
json2line configure formatter -k "new_formatter" -v "{{ .someTemplate }}"

# delete a formatter with the name 'new_formatter'
json2line configure formatter -d "new_formatter"                          
```

Be careful, as this is immediately persisted.

It is also possible to print the current configuration:
```bash
json2line configure -s
```

## Buffer Size

Sometimes a line may be too large and the error message
```
could not parse line: bufio.Scanner: token too long
```
appears. 
For such cases it is possible to set the buffer size to a larger value in the configuration via:
```bash
json2line configure buffer_size -S <BUFFER_SIZE>
```

In case this should not be persisted it is possible to just call
```bash
json2line -b <BUFFER_SIZE>
```
in this case nothing gets persisted and the value in the configuration file gets ignored in
favor of this value.
