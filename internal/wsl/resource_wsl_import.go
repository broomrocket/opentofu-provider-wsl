package wsl

import (
    "context"
    "fmt"
    "unsafe"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "golang.org/x/sys/windows"
)

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lwslapi
// #include "wslapi.h"
// extern long WslRegisterDistribution(PCWSTR distributionName, PCWSTR tarGzFilename);
// extern int WslIsDistributionRegistered(PCWSTR distributionName);
// extern long WslUnregisterDistribution(PCWSTR distributionName);
import "C"

var importDataSource = &schema.Resource{
    CreateContext: importCreate,
    ReadContext:   importRead,
    DeleteContext: importDelete,
    Schema: map[string]*schema.Schema{
        "distribution_name": {
            Type:        schema.TypeString,
            Required:    true,
            ForceNew:    true,
            Description: "Unique name of the distribution.",
        },
        "tar_gz_filename": {
            Type:        schema.TypeString,
            Required:    true,
            ForceNew:    true,
            Description: "Full path to the .tar.gz file that holds the contents of the instance.",
        },
    },
    Description: "The wsl_import resource creates a WSL instance with an image.",
}

func importCreate(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
    distroName, err := windows.UTF16PtrFromString(data.Get("distribution_name").(string))
    if err != nil {
        return diag.Diagnostics{
            diag.Diagnostic{
                Severity:      diag.Error,
                Summary:       "Failed to convert parameter to Windows string",
                Detail:        fmt.Sprintf("%v", err),
                AttributePath: nil,
            },
        }
    }

    tarGZFilename, err := windows.UTF16PtrFromString(data.Get("tar_gz_filename").(string))
    if err != nil {
        return diag.Diagnostics{
            diag.Diagnostic{
                Severity:      diag.Error,
                Summary:       "Failed to convert parameter to Windows string",
                Detail:        fmt.Sprintf("%v", err),
                AttributePath: nil,
            },
        }
    }

    var errors = diag.Diagnostics{}

    var result C.long
    result = C.WslRegisterDistribution((C.PCWSTR)(unsafe.Pointer(distroName)), (C.PCWSTR)(unsafe.Pointer(tarGZFilename)))
    if result == 0x00000000 {
        // S_OK
        data.SetId(data.Get("distribution_name").(string))
    } else {
        // HRESULT_FROM_WIN32(ERROR_ALREADY_EXISTS)
        errors = append(errors, diag.Diagnostic{
            Severity:      diag.Error,
            Summary:       "Distribution failed to register.",
            Detail:        fmt.Sprintf("Error code: 0x%X", uint32(result)),
            AttributePath: nil,
        })
    }
    return errors
}

func importRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
    distroName, err := windows.UTF16PtrFromString(data.Get("distribution_name").(string))
    if err != nil {
        return diag.Diagnostics{
            diag.Diagnostic{
                Severity:      diag.Error,
                Summary:       "Failed to convert parameter to Windows string",
                Detail:        fmt.Sprintf("%v", err),
                AttributePath: nil,
            },
        }
    }
    var result C.int
    result = C.WslIsDistributionRegistered(C.PCWSTR(unsafe.Pointer(distroName)))
    if result == 1 {
        data.SetId(data.Get("distribution_name").(string))
    } else {
        data.SetId("")
    }
    return nil
}

func importDelete(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
    distroName, err := windows.UTF16PtrFromString(data.Get("distribution_name").(string))
    if err != nil {
        return diag.Diagnostics{
            diag.Diagnostic{
                Severity:      diag.Error,
                Summary:       "Failed to convert parameter to Windows string",
                Detail:        fmt.Sprintf("%v", err),
                AttributePath: nil,
            },
        }
    }

    var errors = diag.Diagnostics{}

    var result C.long
    result = C.WslUnregisterDistribution(C.PCWSTR(unsafe.Pointer(distroName)))
    if result == 0x00000000 {
        // S_OK
        data.SetId("")
    } else {
        errors = append(errors, diag.Diagnostic{
            Severity:      diag.Error,
            Summary:       "Distribution failed to unregister.",
            Detail:        fmt.Sprintf("%d", result),
            AttributePath: nil,
        })
    }
    return errors
}
