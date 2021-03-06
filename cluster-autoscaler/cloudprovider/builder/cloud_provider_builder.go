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

package builder

import (
	"io"
	"os"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
<<<<<<< HEAD
	//"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/gce"
=======
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/autoscalr"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/azure"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/gce"
>>>>>>> cluster-autoscaler-release-1.2
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/kubemark"
	"k8s.io/client-go/informers"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	kubemarkcontroller "k8s.io/kubernetes/pkg/kubemark"

	"github.com/golang/glog"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/autoscalr"
)

// AvailableCloudProviders supported by the cloud provider builder.
var AvailableCloudProviders = []string{
	aws.ProviderName,
	azure.ProviderName,
	gce.ProviderNameGCE,
	gce.ProviderNameGKE,
	kubemark.ProviderName,
}

// DefaultCloudProvider is GCE.
const DefaultCloudProvider = gce.ProviderNameGCE

// CloudProviderBuilder builds a cloud provider from all the necessary parameters including the name of a cloud provider e.g. aws, gce
// and the path to a config file
type CloudProviderBuilder struct {
	cloudProviderFlag       string
	cloudConfig             string
	clusterName             string
	autoprovisioningEnabled bool
	regional                bool
}

// NewCloudProviderBuilder builds a new builder from static settings
func NewCloudProviderBuilder(cloudProviderFlag, cloudConfig, clusterName string, autoprovisioningEnabled, regional bool) CloudProviderBuilder {
	return CloudProviderBuilder{
		cloudProviderFlag:       cloudProviderFlag,
		cloudConfig:             cloudConfig,
		clusterName:             clusterName,
		autoprovisioningEnabled: autoprovisioningEnabled,
		regional:                regional,
	}
}

// Build a cloud provider from static settings contained in the builder and dynamic settings passed via args
func (b CloudProviderBuilder) Build(discoveryOpts cloudprovider.NodeGroupDiscoveryOptions, resourceLimiter *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
<<<<<<< HEAD
	var err error
	var cloudProvider cloudprovider.CloudProvider

	nodeGroupsFlag := discoveryOpts.NodeGroupSpecs

	//if b.cloudProviderFlag == "gce" || b.cloudProviderFlag == "gke" {
	//	// GCE Manager
	//	var gceManager gce.GceManager
	//	var gceError error
	//	mode := gce.ModeGCE
	//	if b.cloudProviderFlag == "gke" {
	//		if b.autoprovisioningEnabled {
	//			mode = gce.ModeGKENAP
	//		} else {
	//			mode = gce.ModeGKE
	//		}
	//	}
	//
	//	if b.cloudConfig != "" {
	//		config, fileErr := os.Open(b.cloudConfig)
	//		if fileErr != nil {
	//			glog.Fatalf("Couldn't open cloud provider configuration %s: %#v", b.cloudConfig, err)
	//		}
	//		defer config.Close()
	//		gceManager, gceError = gce.CreateGceManager(config, mode, b.clusterName)
	//	} else {
	//		gceManager, gceError = gce.CreateGceManager(nil, mode, b.clusterName)
	//	}
	//	if gceError != nil {
	//		glog.Fatalf("Failed to create GCE Manager: %v", gceError)
	//	}
	//	cloudProvider, err = gce.BuildGceCloudProvider(gceManager, nodeGroupsFlag, resourceLimiter)
	//	if err != nil {
	//		glog.Fatalf("Failed to create GCE cloud provider: %v", err)
	//	}
	//}

	if b.cloudProviderFlag == "autoscalr" {
		glog.V(0).Infof("provider autoscalr in pass through to aws mode")
		var asrManager *autoscalr.AutoScalrManager
		var awsManager *aws.AwsManager
		var asrError error
		var awsError error
		if b.cloudConfig != "" {
			config, fileErr := os.Open(b.cloudConfig)
			if fileErr != nil {
				glog.Fatalf("Couldn't open cloud provider configuration %s: %#v", b.cloudConfig, err)
			}
			defer config.Close()
			asrManager, asrError = autoscalr.CreateAutoScalrManager(config)
			awsManager, awsError = aws.CreateAwsManager(config)
		} else {
			asrManager, asrError = autoscalr.CreateAutoScalrManager(nil)
			awsManager, awsError = aws.CreateAwsManager(nil)
		}
		if asrError != nil {
			glog.Fatalf("Failed to create AutoScalr Manager: %v", err)
		}
		if awsError != nil {
			glog.Fatalf("Failed to create AWS Manager: %v", err)
		}
		cloudProvider, err = autoscalr.BuildAutoScalrCloudProvider(asrManager, discoveryOpts, resourceLimiter, awsManager)
		if err != nil {
			glog.Fatalf("Failed to create AutoScalr cloud provider: %v", err)
=======
	glog.V(1).Infof("Building %s cloud provider.", b.cloudProviderFlag)
	switch b.cloudProviderFlag {
	case gce.ProviderNameGCE:
		return b.buildGCE(discoveryOpts, resourceLimiter, gce.ModeGCE)
	case gce.ProviderNameGKE:
		if discoveryOpts.DiscoverySpecified() {
			glog.Fatalf("GKE gets nodegroup specification via API, command line specs are not allowed")
		}
		if b.autoprovisioningEnabled {
			return b.buildGCE(discoveryOpts, resourceLimiter, gce.ModeGKENAP)
		}
		return b.buildGCE(discoveryOpts, resourceLimiter, gce.ModeGKE)
	case aws.ProviderName:
		return b.buildAWS(discoveryOpts, resourceLimiter)
	case autoscalr.ProviderName:
		return b.buildAutoScalr(discoveryOpts, resourceLimiter)
	case azure.ProviderName:
		return b.buildAzure(discoveryOpts, resourceLimiter)
	case kubemark.ProviderName:
		return b.buildKubemark(discoveryOpts, resourceLimiter)
	case "":
		// Ideally this would be an error, but several unit tests of the
		// StaticAutoscaler depend on this behaviour.
		glog.Warning("Returning a nil cloud provider")
		return nil
	}

	glog.Fatalf("Unknown cloud provider: %s", b.cloudProviderFlag)
	return nil // This will never happen because the Fatalf will os.Exit
}

func (b CloudProviderBuilder) buildGCE(do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter, mode gce.GcpCloudProviderMode) cloudprovider.CloudProvider {
	var config io.ReadCloser
	if b.cloudConfig != "" {
		var err error
		config, err = os.Open(b.cloudConfig)
		if err != nil {
			glog.Fatalf("Couldn't open cloud provider configuration %s: %#v", b.cloudConfig, err)
>>>>>>> cluster-autoscaler-release-1.2
		}
		defer config.Close()
	}

	manager, err := gce.CreateGceManager(config, mode, b.clusterName, do, b.regional)
	if err != nil {
		glog.Fatalf("Failed to create GCE Manager: %v", err)
	}

	provider, err := gce.BuildGceCloudProvider(manager, rl)
	if err != nil {
		glog.Fatalf("Failed to create GCE cloud provider: %v", err)
	}
	return provider
}

func (b CloudProviderBuilder) buildAWS(do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	var config io.ReadCloser
	if b.cloudConfig != "" {
		var err error
		config, err = os.Open(b.cloudConfig)
		if err != nil {
			glog.Fatalf("Couldn't open cloud provider configuration %s: %#v", b.cloudConfig, err)
		}
		defer config.Close()
	}

	manager, err := aws.CreateAwsManager(config, do)
	if err != nil {
		glog.Fatalf("Failed to create AWS Manager: %v", err)
	}

	provider, err := aws.BuildAwsCloudProvider(manager, rl)
	if err != nil {
		glog.Fatalf("Failed to create AWS cloud provider: %v", err)
	}
	return provider
}

func (b CloudProviderBuilder) buildAutoScalr(do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	var config io.ReadCloser
	if b.cloudConfig != "" {
		var err error
		config, err = os.Open(b.cloudConfig)
		if err != nil {
			glog.Fatalf("Couldn't open cloud provider configuration %s: %#v", b.cloudConfig, err)
		}
		defer config.Close()
	}

	asrManager, asrError := autoscalr.CreateAutoScalrManager(config, do)
	if asrError != nil {
		glog.Fatalf("Failed to create AutoScalr Manager: %v", asrError)
	}
	awsManager, awsError := aws.CreateAwsManager(config, do)
	if awsError != nil {
		glog.Fatalf("Failed to create AWS Manager: %v", awsError)
	}

	provider, err := autoscalr.BuildAutoScalrCloudProvider(asrManager, rl, awsManager)
	if err != nil {
		glog.Fatalf("Failed to create AutoScalr cloud provider: %v", err)
	}
	return provider
}

func (b CloudProviderBuilder) buildAzure(do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	var config io.ReadCloser
	if b.cloudConfig != "" {
		glog.Info("Creating Azure Manager using cloud-config file: %v", b.cloudConfig)
		var err error
		config, err := os.Open(b.cloudConfig)
		if err != nil {
			glog.Fatalf("Couldn't open cloud provider configuration %s: %#v", b.cloudConfig, err)
		}
		defer config.Close()
	} else {
		glog.Info("Creating Azure Manager with default configuration.")
	}
	manager, err := azure.CreateAzureManager(config, do)
	if err != nil {
		glog.Fatalf("Failed to create Azure Manager: %v", err)
	}
	provider, err := azure.BuildAzureCloudProvider(manager, rl)
	if err != nil {
		glog.Fatalf("Failed to create Azure cloud provider: %v", err)
	}
	return provider
}

func (b CloudProviderBuilder) buildKubemark(do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	externalConfig, err := rest.InClusterConfig()
	if err != nil {
		glog.Fatalf("Failed to get kubeclient config for external cluster: %v", err)
	}

	kubemarkConfig, err := clientcmd.BuildConfigFromFlags("", "/kubeconfig/cluster_autoscaler.kubeconfig")
	if err != nil {
		glog.Fatalf("Failed to get kubeclient config for kubemark cluster: %v", err)
	}

	stop := make(chan struct{})

	externalClient := kubeclient.NewForConfigOrDie(externalConfig)
	kubemarkClient := kubeclient.NewForConfigOrDie(kubemarkConfig)

	externalInformerFactory := informers.NewSharedInformerFactory(externalClient, 0)
	kubemarkInformerFactory := informers.NewSharedInformerFactory(kubemarkClient, 0)
	kubemarkNodeInformer := kubemarkInformerFactory.Core().V1().Nodes()
	go kubemarkNodeInformer.Informer().Run(stop)

	kubemarkController, err := kubemarkcontroller.NewKubemarkController(externalClient, externalInformerFactory,
		kubemarkClient, kubemarkNodeInformer)
	if err != nil {
		glog.Fatalf("Failed to create Kubemark cloud provider: %v", err)
	}

	externalInformerFactory.Start(stop)
	if !kubemarkController.WaitForCacheSync(stop) {
		glog.Fatalf("Failed to sync caches for kubemark controller")
	}
	go kubemarkController.Run(stop)

	provider, err := kubemark.BuildKubemarkCloudProvider(kubemarkController, do.NodeGroupSpecs, rl)
	if err != nil {
		glog.Fatalf("Failed to create Kubemark cloud provider: %v", err)
	}
	return provider
}
