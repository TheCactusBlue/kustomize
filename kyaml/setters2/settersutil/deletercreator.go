// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package settersutil

import (
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/openapi"
	"sigs.k8s.io/kustomize/kyaml/setters2"
)

// DeleterCreator delete a definition in the OpenAPI definitions, and removes references
// to the definition from matching resource fields.
type DeleterCreator struct {
	// Name is the name of the setter or substitution to delete
	Name string

	// DefinitionPrefix is the prefix of the OpenAPI definition type
	DefinitionPrefix string

	RecurseSubPackages bool

	OpenAPIFileName string

	// Path to openAPI file
	OpenAPIPath string

	// Path to resources folder
	ResourcesPath string
}

func (d DeleterCreator) Delete() error {
	dd := setters2.DeleterDefinition{
		Name:             d.Name,
		DefinitionPrefix: d.DefinitionPrefix,
	}
	if err := dd.DeleteFromFile(d.OpenAPIPath); err != nil {
		return err
	}

	// Load the updated definitions
	clean, err := openapi.AddSchemaFromFile(d.OpenAPIPath)
	if err != nil {
		return err
	}
	defer clean()

	// Update the resources with the deleter reference
	inout := &kio.LocalPackageReadWriter{PackagePath: d.ResourcesPath, PackageFileName: d.OpenAPIFileName}
	return kio.Pipeline{
		Inputs: []kio.Reader{inout},
		Filters: []kio.Filter{kio.FilterAll(
			&setters2.Delete{
				Name:             d.Name,
				DefinitionPrefix: d.DefinitionPrefix,
			})},
		Outputs: []kio.Writer{inout},
	}.Execute()
}
