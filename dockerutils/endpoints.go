package dockerutils

import "fmt"

func AuthRegistryEndpoint(image string) string {
	return fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/%s:pull", image)
}

func GetManifestEndpoint(image string) string {
	return fmt.Sprintf("https://registry.hub.docker.com/v2/library/%s/manifests/latest", image)
}

func GetBlobFileEndpoint(image, digest string) string {
	return fmt.Sprintf("https://registry.hub.docker.com/v2/library/%s/blobs/%s", image, digest)
}
