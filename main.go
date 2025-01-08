package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/integrations/terraform-provider-github/v6/framework"
	"github.com/integrations/terraform-provider-github/v6/github"
)

func main() {
	ctx := context.Background()
	primary := github.Provider()

	providers := []func() tfprotov5.ProviderServer{
		func() tfprotov5.ProviderServer {
			return schema.NewGRPCProviderServer(primary)
		},
		providerserver.NewProtocol5(framework.New(primary)),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	err = tf5server.Serve("registry.terraform.io/integrations/github", muxServer.ProviderServer, serveOpts...)
}
