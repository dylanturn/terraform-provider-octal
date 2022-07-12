package components

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type octalComponentObject struct {
	group       string
	version     string
	kind        string
	name        string
	namespace   string
	labels      map[string]string
	annotations map[string]string
	manifest    string
	specHash    string
	unstructured.Unstructured
}

type OctalComponentObject interface {
	GetGroup() string
	GetVersion() string
	GetKind() string
	GetName() string
	GetNamespace() string
	GetLabels() map[string]string
	GetAnnotations() map[string]string
	GetManifest() string
	GetSpecHash() string
	GetUnstructuredObject() unstructured.Unstructured
	GetFlat() map[string]interface{}
}

func objectFromManifest(objectManifest string) (OctalComponentObject, error) {

	object, err := contentToObject(objectManifest)
	if err != nil {
		return nil, err
	}

	objectSpec := object.Object["spec"]
	objectSpecJson, marshallErr := json.Marshal(objectSpec)
	if marshallErr != nil {
		return nil, fmt.Errorf("marshall error! %s", marshallErr.Error())
	}

	componentObject := octalComponentObject{
		group:        object.GetObjectKind().GroupVersionKind().Group,
		version:      object.GetObjectKind().GroupVersionKind().Version,
		kind:         object.GetObjectKind().GroupVersionKind().Kind,
		name:         object.GetName(),
		namespace:    object.GetNamespace(),
		labels:       object.GetLabels(),
		annotations:  object.GetLabels(),
		manifest:     objectManifest,
		specHash:     fmt.Sprintf("%x", sha256.Sum256(objectSpecJson)),
		Unstructured: *object,
	}

	return &componentObject, nil
}

func (componentObject *octalComponentObject) GetGroup() string {
	return componentObject.group
}

func (componentObject *octalComponentObject) GetVersion() string {
	return componentObject.version
}

func (componentObject *octalComponentObject) GetKind() string {
	return componentObject.kind
}

func (componentObject *octalComponentObject) GetName() string {
	return componentObject.name
}

func (componentObject *octalComponentObject) GetNamespace() string {
	return componentObject.namespace
}

func (componentObject *octalComponentObject) GetLabels() map[string]string {
	return componentObject.labels
}

func (componentObject *octalComponentObject) GetAnnotations() map[string]string {
	return componentObject.annotations
}

func (componentObject *octalComponentObject) GetManifest() string {
	return componentObject.manifest
}

func (componentObject *octalComponentObject) GetSpecHash() string {
	return componentObject.specHash
}

func (componentObject *octalComponentObject) GetUnstructuredObject() unstructured.Unstructured {
	return componentObject.Unstructured
}

func (componentObject *octalComponentObject) GetFlat() map[string]interface{} {
	return map[string]interface{}{
		"group":       componentObject.GetGroup(),
		"version":     componentObject.GetVersion(),
		"kind":        componentObject.GetKind(),
		"name":        componentObject.GetName(),
		"namespace":   componentObject.GetNamespace(),
		"labels":      componentObject.GetLabels(),
		"annotations": componentObject.GetAnnotations(),
		"spec_hash":   componentObject.GetSpecHash(),
	}
}
