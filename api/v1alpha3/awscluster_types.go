/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

const (
	// ClusterFinalizer allows ReconcileAWSCluster to clean up AWS resources associated with AWSCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "awscluster.infrastructure.cluster.x-k8s.io"
)

// AWSClusterSpec defines the desired state of AWSCluster
type AWSClusterSpec struct {
	// NetworkSpec encapsulates all things related to AWS network.
	NetworkSpec NetworkSpec `json:"networkSpec,omitempty"`

	// The AWS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host.
	SSHKeyName string `json:"sshKeyName,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// ControlPlaneLoadBalancer is optional configuration for customizing control plane behavior
	// +optional
	ControlPlaneLoadBalancer *AWSLoadBalancerSpec `json:"controlPlaneLoadBalancer,omitempty"`

	// ImageLookupOrg is the AWS Organization ID to look up machine images when a
	// machine does not specify an AMI. When set, this will be used for all
	// cluster machines unless a machine specifies a different ImageLookupOrg.
	// +optional
	ImageLookupOrg string `json:"imageLookupOrg,omitempty"`

	// ImageLookupBaseOS is the name of the base operating system used to look
	// up machine images when a machine does not specify an AMI. When set, this
	// will be used for all cluster machines unless a machine specifies a
	// different ImageLookupBaseOS.
	ImageLookupBaseOS string `json:"imageLookupBaseOS,omitempty"`
}

// AWSLoadBalancerSpec defines the desired state of an AWS load balancer
type AWSLoadBalancerSpec struct {

	// Ingress rules associated with the load balancer's Security Group (defaults to
	// 0.0.0.0/0 if not specified)
	// +optional
	IngressCidrBlocks []string `json:ingressCidrBlocks,omitempty"`

	// Scheme sets the scheme of the load balancer (defaults to Internet-facing)
	// +optional
	Scheme *ClassicELBScheme `json:"scheme,omitempty"`
}

// AWSClusterStatus defines the observed state of AWSCluster
type AWSClusterStatus struct {
	Network        Network                  `json:"network,omitempty"`
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`
	Bastion        Instance                 `json:"bastion,omitempty"`
	Ready          bool                     `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AWSCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for EC2 instances"
// +kubebuilder:printcolumn:name="VPC",type="string",JSONPath=".spec.networkSpec.vpc.id",description="AWS VPC the cluster is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.apiEndpoints[0]",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Bastion IP",type="string",JSONPath=".status.bastion.publicIp",description="Bastion IP address for breakglass access"

// AWSCluster is the Schema for the awsclusters API
type AWSCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSClusterSpec   `json:"spec,omitempty"`
	Status AWSClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSClusterList contains a list of AWSCluster
type AWSClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSCluster{}, &AWSClusterList{})
}
