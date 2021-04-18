# tplrender
Go template rendering CLI

## Usage

```sh
go get github.com/bancek/tplrender

tplrender -template - '{"name": "World"}' <<EOF
{{\$data := ((index .Args 0) | fromJson)}}Hello {{\$data.name}}
There is no place like {{.Env.HOME}}
EOF
```

Template context:

- `.Args`: process arguments from `flag.Args()`
- `.Env`: `map[string]string` of the process's environment variables

[sprig](http://masterminds.github.io/sprig/) template functions are available in the template.
