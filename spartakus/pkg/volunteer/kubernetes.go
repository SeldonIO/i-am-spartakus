/*
Copyright 2016 The Kubernetes Authors.

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

package volunteer

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/kubernetes-incubator/spartakus/pkg/report"
	kclient "k8s.io/client-go/1.5/kubernetes"
	kapi "k8s.io/client-go/1.5/pkg/api"
	kv1 "k8s.io/client-go/1.5/pkg/api/v1"
	krest "k8s.io/client-go/1.5/rest"
)

// cloudProviders is a whitelist of the known Kubernetes cloud providers.
var cloudProviders map[string]bool = map[string]bool{
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/aws/aws.go#L55
	"aws": true,
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/azure/azure.go#L34
	"azure": true,
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/cloudstack/cloudstack.go#L30
	"cloudstack": true,
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/gce/gce.go#L55
	"gce": true,
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/mesos/mesos.go#L37
	"mesos": true,
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/openstack/openstack.go#L44
	"openstack": true,
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/ovirt/ovirt.go#L39
	"ovirt": true,
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/photon/photon.go#L43
	"photon": true,
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/rackspace/rackspace.go#L48
	"rackspace": true,
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/vsphere/vsphere.go#L52
	"vsphere": true,
}

type nodeLister interface {
	ListNodes() ([]report.Node, error)
}

type serverVersioner interface {
	ServerVersion() (string, error)
}

func nodeFromKubeNode(kn *kv1.Node) report.Node {
	n := report.Node{
		ID:                      getID(kn),
		OperatingSystem:         strPtr(kn.Status.NodeInfo.OperatingSystem),
		OSImage:                 strPtr(kn.Status.NodeInfo.OSImage),
		KernelVersion:           strPtr(kn.Status.NodeInfo.KernelVersion),
		Architecture:            strPtr(kn.Status.NodeInfo.Architecture),
		ContainerRuntimeVersion: strPtr(kn.Status.NodeInfo.ContainerRuntimeVersion),
		KubeletVersion:          strPtr(kn.Status.NodeInfo.KubeletVersion),
		CloudProvider:           strPtr(providerName(kn.Spec.ProviderID)),
	}
	// We want to iterate the resources in a deterministic order.
	keys := []string{}
	for k, _ := range kn.Status.Capacity {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := kn.Status.Capacity[kv1.ResourceName(k)]
		n.Capacity = append(n.Capacity, report.Resource{
			Resource: string(k),
			Value:    v.String(),
		})
	}
	return n
}

func getID(kn *kv1.Node) string {
	// We don't want to report the node's Name - that is PII.  The MachineID is
	// apparently not always populated and SystemUUID is ill-defined.  Let's
	// just hash them all together.  It should be stable, and this reduces risk
	// of PII leakage.
	return hashOf(kn.Name + kn.Status.NodeInfo.MachineID + kn.Status.NodeInfo.SystemUUID)
}

func hashOf(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil)[0:])
}

func strPtr(str string) *string {
	if str == "" {
		return nil
	}
	p := new(string)
	*p = str
	return p
}

// providerName extracts the cloud provider name from a given
// string that should match: <ProviderName>://<ProviderSpecficNodeID>
// (see https://github.com/kubernetes/client-go/blob/v1.5.1/1.5/pkg/api/v1/types.go#L2446).
// If the given string does not match this format, we return "unknown".
func providerName(providerID string) string {
	parts := strings.Split(providerID, "://")
	if len(parts) == 2 && cloudProviders[parts[0]] {
		return parts[0]
	}
	return "unknown"
}

func newKubeClientWrapper() (*kubeClientWrapper, error) {
	kubeConfig, err := krest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	kubeClient, err := kclient.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return &kubeClientWrapper{client: kubeClient}, nil
}

type kubeClientWrapper struct {
	client *kclient.Clientset
}

func (k *kubeClientWrapper) ListNodes() ([]report.Node, error) {
	knl, err := k.client.Core().Nodes().List(kapi.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodes := make([]report.Node, len(knl.Items))
	for i := range knl.Items {
		kn := &knl.Items[i]
		nodes[i] = nodeFromKubeNode(kn)
	}
	return nodes, nil
}

func (k *kubeClientWrapper) ServerVersion() (string, error) {
	info, err := k.client.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return info.String(), nil
}
