package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-foxycart/foxyprovider"
)

func main() {
	_ = providerserver.Serve(context.Background(), foxyprovider.New, providerserver.ServeOpts{
		// @todo This address will change once the provider is registered in the Terraform Registry
		Address: "67bricks.com/terraform/foxy",
	})
}
