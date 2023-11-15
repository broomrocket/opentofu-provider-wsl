# WSL OpenTofu provider (experimental)

This OpenTofu provider manages Windows Subsystem for Linux instances.

## Compiling

You can compile this provider on Windows by running:

```
go build
```

Note: you must have a cgo/c/c++ toolchain, such as Visual Studio/msbuild or MINGW installed.

## Usage

Currently, the only supported resource is `wsl_import` to import WSL instances:

```hcl2
resource "wsl_import" "test" {
    distribution_name = "test"
    tar_gz_filename = "C:\\path\\to\\distro.tar.gz"
}
```
