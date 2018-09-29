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

package find

import (
	"fmt"

	bpb "github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/apis/bundle/v1alpha1"
	"github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/converter"
	"github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/core"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

// BundleFinder is a wrapper which allows for efficient searching through
// bundles. The BundleFinder is intended to be readonly; if modifications are
// made to the bundle, subsequent lookups will fail.
type BundleFinder struct {
	bundle        *bpb.ClusterBundle
	nodeLookup    map[string]*bpb.NodeConfig
	compLookup    map[string]*bpb.ClusterComponent
	compObjLookup map[core.ClusterObjectKey][]*structpb.Struct
}

// NewBundleFinder creates a new BundleFinder or returns an error.
func NewBundleFinder(b *bpb.ClusterBundle) (*BundleFinder, error) {
	b = converter.CloneBundle(b)
	nodeConfigs := make(map[string]*bpb.NodeConfig)
	for _, nc := range b.GetSpec().GetNodeConfigs() {
		n := nc.GetMetadata().GetName()
		if n == "" {
			return nil, fmt.Errorf("node bootstrap configs must always have a metadata.name. was empty for %v", nc)
		}
		nodeConfigs[n] = nc
	}

	compConfigs := make(map[string]*bpb.ClusterComponent)
	for _, ca := range b.GetSpec().GetComponents() {
		n := ca.GetMetadata().GetName()
		if n == "" {
			return nil, fmt.Errorf("cluster components must always have a metadata.name. was empty for %v", ca)
		}
		compConfigs[n] = ca
	}

	return &BundleFinder{
		bundle:     b,
		nodeLookup: nodeConfigs,
		compLookup: compConfigs,
	}, nil
}

// ClusterComponent returns a found cluster component or nil.
func (b *BundleFinder) ClusterComponent(name string) *bpb.ClusterComponent {
	return b.compLookup[name]
}

// NodeConfig returns a node bootstrap config or nil.
func (b *BundleFinder) NodeConfig(name string) *bpb.NodeConfig {
	return b.nodeLookup[name]
}

// ClusterComponentObject returns ClusterComponent's Cluster object or nil.
func (b *BundleFinder) ClusterComponentObject(compName string, objName string) []*structpb.Struct {
	comp := b.ClusterComponent(compName)
	var out []*structpb.Struct
	for _, c := range comp.GetClusterObjects() {
		n := core.ObjectName(c)
		if n == objName {
			out = append(out, c)
		}
	}
	return out
}