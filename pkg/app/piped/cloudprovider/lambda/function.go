// Copyright 2021 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lambda

import (
	"fmt"
	"io/ioutil"
	"strings"

	"sigs.k8s.io/yaml"
)

const (
	versionV1Beta1       = "pipecd.dev/v1beta1"
	functionManifestKind = "LambdaFunction"
	// Memory and Timeout lower and upper limit as noted via
	// https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/#UpdateFunctionConfigurationInput
	memoryLowerLimit  = 1
	timeoutLowerLimit = 1
	timeoutUpperLimit = 900
)

type FunctionManifest struct {
	Kind       string               `json:"kind"`
	APIVersion string               `json:"apiVersion,omitempty"`
	Spec       FunctionManifestSpec `json:"spec"`
}

func (fm *FunctionManifest) validate() error {
	if fm.APIVersion != versionV1Beta1 {
		return fmt.Errorf("unsupported version: %s", fm.APIVersion)
	}
	if fm.Kind != functionManifestKind {
		return fmt.Errorf("invalid manifest kind given: %s", fm.Kind)
	}
	if err := fm.Spec.validate(); err != nil {
		return err
	}
	return nil
}

// FunctionManifestSpec contains configuration for LambdaFunction.
type FunctionManifestSpec struct {
	Name            string            `json:"name"`
	Role            string            `json:"role"`
	ImageURI        string            `json:"image"`
	S3Bucket        string            `json:"s3Bucket"`
	S3Key           string            `json:"s3Key"`
	S3ObjectVersion string            `json:"s3ObjectVersion"`
	Handler         string            `json:"handler"`
	Memory          int32             `json:"memory"`
	Timeout         int32             `json:"timeout"`
	Tags            map[string]string `json:"tags,omitempty"`
	Environments    map[string]string `json:"environments,omitempty"`
}

func (fmp FunctionManifestSpec) validate() error {
	if len(fmp.Name) == 0 {
		return fmt.Errorf("lambda function is missing")
	}
	if len(fmp.ImageURI) == 0 && len(fmp.S3Bucket) == 0 {
		return fmt.Errorf("one of image or s3 bucket is required to be configured")
	}
	if len(fmp.Role) == 0 {
		return fmt.Errorf("role is missing")
	}
	if fmp.Memory < memoryLowerLimit {
		return fmt.Errorf("memory is missing")
	}
	if fmp.Timeout < timeoutLowerLimit || fmp.Timeout > timeoutUpperLimit {
		return fmt.Errorf("timeout is missing or out of range")
	}
	return nil
}

func loadFunctionManifest(path string) (FunctionManifest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return FunctionManifest{}, err
	}
	return parseFunctionManifest(data)
}

func parseFunctionManifest(data []byte) (FunctionManifest, error) {
	var obj FunctionManifest
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return FunctionManifest{}, err
	}
	if err := obj.validate(); err != nil {
		return FunctionManifest{}, err
	}
	return obj, nil
}

// DecideRevisionName returns revision name to apply.
func DecideRevisionName(fm FunctionManifest, commit string) (string, error) {
	tag, err := FindImageTag(fm)
	if err != nil {
		return "", err
	}
	tag = strings.ReplaceAll(tag, ".", "")

	if len(commit) > 7 {
		commit = commit[:7]
	}
	return fmt.Sprintf("%s-%s-%s", fm.Spec.Name, tag, commit), nil
}

// FindImageTag parses image tag from given LambdaFunction manifest.
func FindImageTag(fm FunctionManifest) (string, error) {
	name, tag := parseContainerImage(fm.Spec.ImageURI)
	if name == "" {
		return "", fmt.Errorf("image name could not be empty")
	}
	return tag, nil
}

func parseContainerImage(image string) (name, tag string) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		tag = parts[1]
	}
	paths := strings.Split(parts[0], "/")
	name = paths[len(paths)-1]
	return
}
