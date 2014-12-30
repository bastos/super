# Super

Super is a simple, WIP command line tool to launch programs based on user created patterns.

You can use Super to open URLs, run scripts or compile programs.

Super is based on tools like [Plumber](http://en.wikipedia.org/wiki/Plumber_%28program%29) and [Alfred](http://www.alfredapp.com/).

## Instalation

```go get github.com/bastos/super```

Create a TOML configuration file located at ```~/.super.toml```.

## Example

Rule:

```toml
[[rule]]
name = "Jira"
regex = "^SU\\-([0-9]*)$"
command = "open https://COMPANY.atlassian.net/browse/$1"

[[rule]]
name = "Github"
regex = "^(.*)/(.*)$"
prefix = "gh:"
command = "open https://github.com/$1"
```

Running:

```super SB-3495```

```super gh:bastos/super```

## Verify you configuration

```super --check```
