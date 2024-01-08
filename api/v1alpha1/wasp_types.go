/*
Copyright 2024.

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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SwapStrategy string

const (
	Orthogonal SwapStrategy = "orthogonal"
	AllowSpike SwapStrategy = "allow-spike"
	forceHard  SwapStrategy = "force-hard"
)

// WaspSpec defines the desired state of Wasp
type WaspSpec struct {
	// Defines the reclaim strategy with regards to swapping
	Strategy SwapStrategy `json:"strategy,omitempty"`
	// The size of the swap file
	SwapFileSize *resource.Quantity `json:"swapFileSize,omitempty"`
	// The name of the swap file
	SwapFileName string `json:"swapFileName,omitempty"`
	// Path where the swap file should be created
	SwapFilePath string `json:"SwapFilePath,omitempty"`
	// The chroot path of the host root filesystem
	FsRoot string `json:"fsRoot,omitempty"`
	// level of chatiness for debug purpose
	verbosity int `json:"verbosity,omitempty"`
}

type WaspConditionType string

const (
	SwapConfigurtaionConditionDeployed         WaspConditionType = "Deployed"
	SwapConfigurtaionConditionDeployInProgress WaspConditionType = "DeployInProgress"
	SwapConfigurtaionConditionFailed           WaspConditionType = "Failed"
)

type WaspCondition struct {
	Type   WaspConditionType      `json:"type"`
	Status metav1.ConditionStatus `json:"status"`
	// +optional
	// +nullable
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// +optional
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

// WaspStatus defines the observed state of Wasp
type WaspStatus struct {
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []WaspCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Wasp is the Schema for the wasps API
type Wasp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WaspSpec   `json:"spec,omitempty"`
	Status WaspStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WaspList contains a list of Wasp
type WaspList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Wasp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Wasp{}, &WaspList{})
}
