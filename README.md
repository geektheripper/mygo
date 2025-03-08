# MyGo

A Lifetime Golang Monorepo Manager.

## Directory Structure Example

```
.
...
├── pkg
│   └── nagasaki
│   └── soyo
├── README.md
├── LICENSE
...
```

mygo publish sub packages with only necessary files.

for `pkg/nagasaki` in example project structure, mygo publish following files to sub package tag:

```
pkg/nagasaki/...
LICENSE
```

## Usage

```bash
# create a new sub package (just create directory)
mygo new pkg/tomori

# publish a sub package, auto find latest tag:
#    upgrade semver patch part, if package
mygo publish pkg/tomori

# publish a sub package, auto find latest tag, upgrade semver and push
mygo publish pkg/tomori --patch
mygo publish pkg/tomori --minor
mygo publish pkg/tomori --major

# publish with exact version (without "v" prefix)
mygo publish pkg/tomori --version 1.2.3

# add more files to publish, known issue:
#   - can't use glob pattern, only support exact file name
#   - only support file in project root dir
mygo publish pkg/tomori --includes "README.md" --includes "another file"

# not publish license file
mygo publish pkg/tomori --no-license

# publish sub package with message
mygo publish pkg/tomori --message "add new feature"
```
