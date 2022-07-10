package util

import (
	"bytes"
	"context"
	"embed"
	"fmt"

	Appsv1 "k8s.io/api/apps/v1"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
)

func LoadDeploymentObj(ctx context.Context, objManifest string) *Appsv1.Deployment {
	deploymentObject := &Appsv1.Deployment{}
	LoadManifest(objManifest).Decode(&deploymentObject)
	return deploymentObject
}

func DecodeManifest(manifest []byte) *k8Yaml.YAMLOrJSONDecoder {
	return k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(manifest), 1000)
}

func LoadManifest(manifestPath string) *k8Yaml.YAMLOrJSONDecoder {
	return k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(readFile(manifestPath)), 1000)
}

func ReadEmbeddedFiles(embeddedFs embed.FS) []string {
	embeddedDirectories, err := embeddedFs.ReadDir(".")
	if err != nil {
		fmt.Println("Failed to read embedded directory", err.Error())
	}

	for _, directory := range embeddedDirectories {

		fileList, err := embed.FS.ReadDir(embeddedFs, directory.Name())
		if err != nil {
			fmt.Println("Failed to read embedded directory", err.Error())
		}

		fileContentStrings := make([]string, len(fileList))

		for index, file := range fileList {
			fileContents, err := embeddedFs.ReadFile(fmt.Sprintf("%s/%s", directory.Name(), file.Name()))
			if err != nil {
				fmt.Println("Failed to read file", err.Error())
			}
			fileContentStrings[index] = string(fileContents)
		}
		return fileContentStrings
	}
	return []string{}
}
