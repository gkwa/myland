````
-- README.md --
# myland

A tool to escape template delimiters in files, primarily for use with [boilerplate](https://github.com/gruntwork-io/boilerplate).

## Problem

When using [boilerplate](https://github.com/gruntwork-io/boilerplate) with files that contain template delimiters (e.g. `{{}}` in justfile), there is a conflict since boilerplate uses these same delimiters for its own templating.

For example, this justfile line:

```makefile
test -z "{{ shell_files }}" || shfmt -w -s -i 4 {{ shell_files }}
````

Conflicts with boilerplate's template processing. The solution is to escape the delimiters:

```makefile
test -z "{{"{{"}} shell_files {{"}}"}}" || shfmt -w -s -i 4 {{"{{"}} shell_files {{"}}"}}
```

While boilerplate has a proposed `skip_templating` feature ([PR #184](https://github.com/gruntwork-io/boilerplate/pull/184)), it's not yet merged. This tool provides an automated way to escape delimiters in the meantime.

## Usage

To escape delimiters in a file:

```bash
myland justfile
```

This will process the file in place, converting:

- `{{` to `{{"{{"}}`
- `}}` to `{{"}}"}}`

## Install

```bash
go install github.com/gkwa/myland@latest
```

## Example

Input justfile:

```makefile
set shell := ["bash", "-uec"]
shell_files := `find . -name .git -prune -o -name '*.sh' -print`
fmt:
    test -z "{{ shell_files }}" || shfmt -w -s -i 4 {{ shell_files }}
```

After running `myland justfile`:

```makefile
set shell := ["bash", "-uec"]
shell_files := `find . -name .git -prune -o -name '*.sh' -print`
fmt:
    test -z "{{"{{"}} shell_files {{"}}"}}" || shfmt -w -s -i 4 {{"{{"}} shell_files {{"}}"}}
```

The escaped file can now be used with boilerplate without conflicts.

## Development

Requires:

- Go 1.21+
- make

Setup:

```bash
git clone https://github.com/gkwa/myland
cd myland
make build
```
