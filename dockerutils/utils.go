package dockerutils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/codecrafters-io/docker-starter-go/throwerror"
)

var mediaType string = "application/vnd.docker.distribution.manifest.v2+json"

type AuthResponse struct {
	Token string `json:"token"`
}

type manifestLayer struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int    `json:"size"`
}

type Manifest struct {
	SchemaVersion int             `json:"schemaVersion"`
	MediaType     string          `json:"mediaType"`
	Layers        []manifestLayer `json:"layers"`
}

func GetAuthToken(image string) *AuthResponse {
	res, err := http.Get(AuthRegistryEndpoint(image))
	if err != nil {
		throwerror.ThrowError(err, "Unable to get auth-token from docker")
	}

	var authToken AuthResponse
	err = json.NewDecoder(res.Body).Decode(&authToken)
	if err != nil {
		throwerror.ThrowError(err, "Error encoding get-auth-token res-body into AuthResponse type")
	}

	return &authToken
}

func GetManifest(image string, authToken *AuthResponse) Manifest {
	req, err := http.NewRequest("GET", GetManifestEndpoint(image), nil)
	if err != nil {
		throwerror.ThrowError(err, "Unable to create request to get manifest")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken.Token))
	req.Header.Set("Accept", mediaType)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		throwerror.ThrowError(err, "Unable to get manifest from registry")
	}

	var manifest Manifest
	err = json.NewDecoder(res.Body).Decode(&manifest)
	if err != nil {
		throwerror.ThrowError(err, "Unable to decode get-manifest res-body into Manifest type")
	}

	if manifest.MediaType != mediaType {
		// lets use the previous value of err itself
		throwerror.ThrowError(err, "Mediatype in manifest response not matching with mediaType in request (Accept header)")
	}

	return manifest
}

func DownloadAndExtractLayers(layers []manifestLayer, image string, authToken *AuthResponse, tempDirPath string) {
	for _, layer := range layers {
		req, err := http.NewRequest("GET", GetBlobFileEndpoint(image, layer.Digest), nil)
		if err != nil {
			throwerror.ThrowError(err, "Unable to create request to download layer")
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken.Token))

		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			throwerror.ThrowError(err, "Unable to download blob from registry")
		}

		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			throwerror.ThrowError(err, "Could not read response body downloaded from blob registry")
		}

		filename := "image.tar"
		file, err := os.OpenFile(filename, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			throwerror.ThrowError(err, "Unable to open new file to write blob to")
		}

		_, err = file.Write(resBody)
		if err != nil {
			throwerror.ThrowError(err, "Unable to write blob to newly created file")
		}

		exec.Command("tar", "-xf", filename, "-C", tempDirPath).Run()
	}
}
