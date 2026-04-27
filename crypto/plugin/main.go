package main

import (
	"github.com/tengolang/tengo-modules/crypto"
	"github.com/tengolang/tengo/v3"
)

// TengoModule is the symbol loaded by tengo's PluginLoader.
var TengoModule map[string]tengo.Object = crypto.Module

func main() {}
