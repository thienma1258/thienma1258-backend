package main
import (
	artifactregistry "cloud.google.com/go/artifactregistry/apiv1beta2"
	"context"
)
func client() {
	ctx := context.Background()
	c, err := artifactregistry.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

}
