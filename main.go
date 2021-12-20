package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8yaml "k8s.io/apimachinery/pkg/util/yaml"
)

func main() {
	rawdata, err := ioutil.ReadFile("configMap.yaml")
	if err != nil {
		panic(err)
	}
	decoder := k8yaml.NewYAMLOrJSONDecoder(bytes.NewReader(rawdata), 1)
	var manifests unstructured.Unstructured
	for {
		nxtManifest := unstructured.Unstructured{}
		err := decoder.Decode(&nxtManifest)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		// Skip empty manifests
		if len(nxtManifest.Object) > 0 {
			manifests = nxtManifest
		}
	}
	fmt.Println(manifests)
	replaceInner(&manifests.Object, configReplacement)
	fmt.Println(manifests)
}
