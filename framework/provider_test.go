package framework

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/integrations/terraform-provider-github/v6/github"
)

func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"github": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			primary := github.Provider()

			providers := []func() tfprotov5.ProviderServer{
				func() tfprotov5.ProviderServer {
					return schema.NewGRPCProviderServer(primary)
				},
				providerserver.NewProtocol5(New(primary)),
			}

			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}

func TestMuxServer(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `provider "github" {}`,
			},
		},
	})
}
