## Getting Started

### Install
```sh
go get github.com/hmarui66/go-watch
```

### Start watch
```sh
cd [project root]
go-watch
```

#### Go build flags

You can use go build flags with `-build-flags` flag.
These flags are set in the form of `go build -o {binfile} [build flags]`.

For example, this command

```sh
go-watch -build-flags '-ldflags "-X configs.someValue=100 -X configs.anotherValue=300"'
```

will trigger the following build command.

```sh
go build -o {binfile} -ldflags "-X configs.someValue=100 -X configs.anotherValue=300"
```

## Refs

https://github.com/pilu/fresh
