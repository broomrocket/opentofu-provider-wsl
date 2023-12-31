package main

import (
    "flag"

    "github.com/broomrocket/opentofu-provider-wsl/internal/wsl"
    "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
    var debugMode bool

    flag.BoolVar(
        &debugMode,
        "debug",
        false,
        "set to true to run the provider with support for debuggers like delve",
    )
    flag.Parse()

    opts := &plugin.ServeOpts{
        ProviderAddr: "registry.terraform.io/broomrocket/wsl",
        ProviderFunc: wsl.New,
    }

    if debugMode {
        opts.Debug = true
    }

    plugin.Serve(opts)
}
