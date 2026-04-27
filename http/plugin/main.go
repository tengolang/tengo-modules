package main

import (
	"github.com/tengolang/tengo-modules/http"
	"github.com/tengolang/tengo/v3"
)

// TengoModule is the symbol loaded by tengo's PluginLoader.
var TengoModule map[string]tengo.Object = http.Module

func main() {}
