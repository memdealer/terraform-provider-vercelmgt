package main

//Garbage-In, Garbage-Out (GIGO)

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
	"vercel-mgt/vercel"
)

func main() {
	err := providerserver.Serve(context.Background(), vercel.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/vercel/vercel",
	})
	if err != nil {
		log.Fatalf("unable to serve provider: %s", err)
	}

}
