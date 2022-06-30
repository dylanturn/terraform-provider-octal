package octal

import (
	"bytes"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
	"log"
)

func readFile(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func loadManifest(manifestPath string) *k8Yaml.YAMLOrJSONDecoder {
	return k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(readFile(manifestPath)), 1000)
}

func octalListOptions(resourceId string) metav1.ListOptions {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"project-octal.io/cert-manager": resourceId}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	return listOptions
}

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}
