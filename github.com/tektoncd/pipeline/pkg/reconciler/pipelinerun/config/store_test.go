/*
Copyright 2019 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline"
	"github.com/tektoncd/pipeline/pkg/artifacts"
	ttesting "github.com/tektoncd/pipeline/pkg/reconciler/testing"
	"github.com/tektoncd/pipeline/test/diff"

	test "github.com/tektoncd/pipeline/pkg/reconciler/testing"
)

func TestStoreLoadWithContext(t *testing.T) {
	store := NewStore(pipeline.Images{}, ttesting.TestLogger(t))
	bucketConfig := test.ConfigMapFromTestFile(t, "config-artifact-bucket")
	store.OnConfigChanged(bucketConfig)

	config := FromContext(store.ToContext(context.Background()))

	expected, _ := artifacts.NewArtifactBucketConfigFromConfigMap(pipeline.Images{})(bucketConfig)
	if d := cmp.Diff(expected, config.ArtifactBucket); d != "" {
		t.Errorf("Unexpected controller config %s", diff.PrintWantGot(d))
	}
}
func TestStoreImmutableConfig(t *testing.T) {
	store := NewStore(pipeline.Images{}, ttesting.TestLogger(t))
	store.OnConfigChanged(test.ConfigMapFromTestFile(t, "config-artifact-bucket"))

	config := store.Load()

	config.ArtifactBucket.Location = "mutated"

	newConfig := store.Load()

	if newConfig.ArtifactBucket.Location == "mutated" {
		t.Error("Controller config is not immutable")
	}
}
