# tengo-modules

Community modules for the [Tengo](https://github.com/tengolang/tengo) scripting language,
extending the standard library with functionality not suitable for the core.

## Usage

Each module lives in its own subdirectory and is importable as a Go package that
registers itself as a Tengo module.

```go
import (
    "github.com/tengolang/tengo/v3"
    // example: "github.com/tengolang/tengo-modules/uuid"
)

modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
// modules.AddBuiltinModule("uuid", uuid.Module)

script := tengo.NewScript(src)
script.SetImports(modules)
```

## Modules

| Module | Description |
| :----- | :---------- |
| _(none yet)_ | |

## Contributing

Each module is a Go package that exports a `Module` variable of type `*tengo.BuiltinModule`.
Add a subdirectory, implement the module, and register it in this README.
