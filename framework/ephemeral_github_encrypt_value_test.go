package framework

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccEphemeralResourceEncryptValue(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		// Ephemeral resources are only available in 1.10 and later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},

		Steps: []resource.TestStep{
			{
				Config: `
				provider "echo" {
					data = ephemeral.github_encrypt_value.test
				}
				resource "echo" "test" {}

				provider "github"{
					owner = "SBGoods"
				}

				data "github_actions_public_key" "test" {
					repository = "terraform-provider-releasetest"
				}

				ephemeral "github_encrypt_value" "test" {
					public_encrypt_key = data.github_actions_public_key.test.key
					plaintext_value = "My_Secret_Value"
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("plaintext_value"),
						knownvalue.StringExact("My_Secret_Value")),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("public_encrypt_key"),
						knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("encrypted_value"),
						knownvalue.NotNull()),
				},
			},
		},
	})
}
