package wsl

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path"
    "path/filepath"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const alpineDownloadLink = "https://dl-cdn.alpinelinux.org/alpine/v3.18/releases/x86_64/alpine-minirootfs-3.18.4-x86_64.tar.gz"

func TestImport(t *testing.T) {
    tempDir := t.TempDir()

    resp, err := http.Get(alpineDownloadLink)
    if err != nil {
        t.Skipf("Failed to download Alpine test image (%v)", err)
        return
    }
    defer func() {
        _ = resp.Body.Close()
    }()
    if resp.StatusCode != 200 {
        t.Skipf("Incorrect status code for Alpine test image (%d)", resp.StatusCode)
        return
    }
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Skipf("Failed to download Alpine test image (%v)", err)
        return
    }
    testImagePath := path.Join(tempDir, "test.tar.gz")
    if err := os.WriteFile(testImagePath, data, 0777); err != nil {
        t.Skipf("Failed to write Alpine test image to disk (%v)", err)
        return
    }

    config := fmt.Sprintf(`
provider "wsl" {}

resource "wsl_import" "test" {
    distribution_name = "wsl_tf_test"
    tar_gz_filename = "%s"
}
`, filepath.ToSlash(testImagePath))

    resource.UnitTest(t, resource.TestCase{
        ProviderFactories: map[string]func() (*schema.Provider, error){
            "wsl": func() (*schema.Provider, error) {
                return New(), nil
            },
        },
        Steps: []resource.TestStep{
            {
                Config: config,
                Check: func(state *terraform.State) error {
                    if state.RootModule().Resources["wsl_import.test"].Primary.ID != "wsl_tf_test" {
                        return fmt.Errorf(
                            "incorrect ID after run: %s",
                            state.RootModule().Resources["wsl_import.test"].Primary.ID,
                        )
                    }
                    return nil
                },
            },
        },
    })
}
