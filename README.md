# Tastevin

A collection of package and tools for testing LinuxBoot


## Gotest

The `gotest` command is a wrapper for running `go test --json` and preparing test data for rendering them with [GoTestWeb](https://github.com/JulienVdG/gotestweb)

### Usage

#### Run tests
```sh
gotest run go test --json ./...
```

#### View results
```sh
gotest serve
```

#### View tests running
```sh
gotest live go test --json ./...
```

#### Generate data for external web-server (CI)
```sh
gotest gen
```

### build a standalone binary (including resources)

This is possible using [go.rice](https://github.com/GeertJohan/go.rice), first install the `rice` tool.
Use one of the following command set to build:

1. by generating go code
```sh
rice embed-go -i ./pkg/gotestweb/
go build ./cmds/gotest/
```

2. by appending an archive
```sh
go build ./cmds/gotest/
rice -i ./pkg/gotestweb/ append --exec gotest
```

In both case you get a `gotest` executable that embed the GoTestWeb data.
You can move it anywhere and still can use `gotest gen`.

