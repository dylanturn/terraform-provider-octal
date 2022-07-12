package components

import (
	"embed"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func contentToObject(content string) (*unstructured.Unstructured, error) {
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(content), 4096)

	var object *unstructured.Unstructured

	for {
		err := decoder.Decode(&object)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal manifest: %s", err)
		}

		if object != nil {
			return object, nil
		}
	}
}

func flattenConfig(prefix string, value interface{}, flatmap map[string]string) {
	submap, ok := value.(map[interface{}]interface{})
	if ok {
		for k, v := range submap {
			flattenConfig(prefix+"."+k.(string), v, flatmap)
		}
		return
	}
	stringlist, ok := value.([]interface{})
	if ok {
		flattenConfig(fmt.Sprintf("%s.size", prefix), len(stringlist), flatmap)
		for i, v := range stringlist {
			flattenConfig(fmt.Sprintf("%s.%d", prefix, i), v, flatmap)
		}
		return
	}
	flatmap[prefix] = fmt.Sprintf("%v", value)
}

func ReadEmbeddedFiles(embeddedFs embed.FS) []string {
	embeddedFiles := []string{}
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
		embeddedFiles = append(embeddedFiles, fileContentStrings...)
	}
	return embeddedFiles
}
