// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transformer

import (
	"fmt"

	bpb "github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/apis/bundle/v1alpha1"
	"github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/find"
)

// ComponentExporter is a struct that exports cluster components.
type ComponentExporter struct {
	finder *find.BundleFinder
}

// NewComponentExporter creates a new component exporter.
func NewComponentExporter(b *bpb.ClusterBundle) (*ComponentExporter, error) {
	f, err := find.NewBundleFinder(b)
	if err != nil {
		return nil, err
	}
	return &ComponentExporter{
		finder: f,
	}, nil
}

// Export extracts the named ClusterComponent from the given bundle. It returns a list of
// ExportedComponents.
// - Returns an error if no component by the given compName is found.
// - Returns an error if the desired component in the given bundle is not inlined.
func (e *ComponentExporter) Export(compName string) (*bpb.ClusterComponent, error) {
	comp := e.finder.ClusterComponent(compName)
	if comp == nil {
		return nil, fmt.Errorf("could not find cluster component named %q", compName)
	}
	return comp, nil
}