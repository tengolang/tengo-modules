// Package main is the plugin entry point for the uuid module.
//
// Build with:
//
//	go build -buildmode=plugin -o uuid.so ./uuid/plugin/
//
// Then place uuid.so in a directory on TENGO_MODULE_PATH or ~/.tengo/modules/
// and import it from Tengo scripts with:
//
//	uuid := import("uuid")
package main

import (
	"github.com/tengolang/tengo/v3"
	"github.com/tengolang/tengo-modules/uuid"
)

// TengoModule is the plugin entry point recognised by tengo's PluginLoader.
var TengoModule map[string]tengo.Object = uuid.Module

func main() {}
