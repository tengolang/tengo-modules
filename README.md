# tengo-modules

Community modules for the [Tengo](https://github.com/tengolang/tengo) scripting language,
extending the standard library with functionality not suitable for the core.

## Module types

**Source modules** (`.tengo` files) work on all platforms and require no compilation.
Drop them in a directory on `TENGO_MODULE_PATH` or `~/.tengo/modules/` and they are
importable immediately.

**Plugin modules** (Go, compiled to `.so`) give scripts access to the full Go ecosystem.
They must be built with the same Go toolchain and the same version of
`github.com/tengolang/tengo/v3` as the running `tengo` binary. Plugins are not
supported on Windows.

## Usage

### Source modules

```sh
# point tengo at the directory containing your .tengo module files
export TENGO_MODULE_PATH=/path/to/tengo-modules/set:$TENGO_MODULE_PATH
```

```tengo
set := import("set")
set.union([1,2,3], [3,4,5])   // [1 2 3 4 5]
```

### Plugin modules

Build the plugin against the same `tengo` version the binary was built with:

```sh
# from the tengo-modules repo root
go build -buildmode=plugin -o ~/.tengo/modules/uuid.so ./uuid/plugin/
```

```tengo
uuid := import("uuid")
uuid.v4()   // "3efe329b-28e8-4ed1-9dc5-b2b7a9a5852a"
```

### Embedding (Go)

```go
import (
    "github.com/tengolang/tengo/v3/stdlib"
    "github.com/tengolang/tengo-modules/uuid"
)

modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
modules.AddBuiltinModule("uuid", uuid.Module)
```

## Modules

| Module | Type | Description |
| :----- | :--- | :---------- |
| [set](set/set.tengo) | source | Set operations on arrays: `from`, `contains`, `union`, `intersect`, `difference`, `equal` |
| [uuid](uuid/uuid.go) | plugin / embed | UUID generation and parsing: `v4`, `v1`, `parse`, `valid`, `nil` |

## Contributing

- **Source module**: add a `<name>/<name>.tengo` file that `export`s a map of functions.
- **Plugin module**: add a `<name>/` Go package (for embedding) and a `<name>/plugin/main.go`
  that exports `var TengoModule map[string]tengo.Object`.
- Register the module in this README.
