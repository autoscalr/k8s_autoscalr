/*
Copyright 2017 AutoScalr

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

package autoscalr

import (
<<<<<<< HEAD
	"os"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/kubernetes/plugin/pkg/scheduler/schedulercache"
	"github.com/golang/glog"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
=======
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/kubernetes/pkg/scheduler/schedulercache"
	"k8s.io/apimachinery/pkg/api/resource"
	"github.com/golang/glog"
	//"k8s.io/client-go/rest"
	kube_client "k8s.io/client-go/kubernetes"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"net/url"
	"os"
	"fmt"
	"time"
)

const (
	// ProviderName is the cloud provider name for AWS
	ProviderName = "autoscalr"
>>>>>>> cluster-autoscaler-release-1.2
)

// autoScalrCloudProvider implements CloudProvider interface.
type autoScalrCloudProvider struct {
	autoScalrManager	*AutoScalrManager
	awsProvider			cloudprovider.CloudProvider
}

<<<<<<< HEAD
func BuildAutoScalrCloudProvider(autoScalrManager *AutoScalrManager, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions, resourceLimiter *cloudprovider.ResourceLimiter, awsManager *aws.AwsManager) (*autoScalrCloudProvider, error) {
	awsProv, err := aws.BuildAwsCloudProvider(awsManager, discoveryOpts, resourceLimiter)
	if err != nil {
		glog.V(0).Infof("Received error from BuildAwsCloudProvider: %s", err.Error())
	} else {
		for _, spec := range discoveryOpts.NodeGroupSpecs {
			specObj, err := dynamic.SpecFromString(spec, true)
			if err != nil {
				glog.V(0).Infof("Received error from SpecFromString: %s", err.Error())
			} else {
				// Record ASG in env variable
				os.Setenv("AUTOSCALING_GROUP_NAME", specObj.Name)
			}
		}
=======
func BuildAutoScalrCloudProvider(autoScalrManager *AutoScalrManager, resourceLimiter *cloudprovider.ResourceLimiter, awsManager *aws.AwsManager) (*autoScalrCloudProvider, error) {
	awsProv, err := aws.BuildAwsCloudProvider(awsManager, resourceLimiter)
	if err != nil {
		glog.V(0).Infof("Received error from BuildAwsCloudProvider: %s", err.Error())
	} else {
>>>>>>> cluster-autoscaler-release-1.2
		createAsrAppIfNeeded()
		provider := &autoScalrCloudProvider{
			autoScalrManager: autoScalrManager,
			awsProvider:      awsProv,
		}
		return provider, err
	}
	return nil, err
}

func createAsrAppIfNeeded() error {
	err := appDefCreate()
	return err
}

<<<<<<< HEAD
=======
// Cleanup stops the go routine that is handling the current view of the ASGs in the form of a cache
func (asrProvider *autoScalrCloudProvider) Cleanup() error {
	return nil
}

>>>>>>> cluster-autoscaler-release-1.2
// Name returns name of the cloud provider.
func (asrProvider *autoScalrCloudProvider) Name() string {
	return "autoscalr"
}

// NodeGroups returns all node groups configured for this cloud provider.
func (asrProvider *autoScalrCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	awsNGs := asrProvider.awsProvider.NodeGroups()
	asrNGs := make([]cloudprovider.NodeGroup, 0, len(awsNGs))
	for _, nodeGrp := range awsNGs {
		asrNGs = append(asrNGs, BuildAutoScalrNodeGroup(nodeGrp))
	}
	return asrNGs
}

// NodeGroupForNode returns the node group for the given node.
func (asrProvider *autoScalrCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	awsNg, err := asrProvider.awsProvider.NodeGroupForNode(node)
	if err != nil {
		return awsNg, err
	} else {
		// wrap in asrNode
		return BuildAutoScalrNodeGroup(awsNg), err
	}
}

// Pricing returns pricing model for this cloud provider or error if not available.
func (asrProvider *autoScalrCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	return asrProvider.awsProvider.Pricing()
}

// GetAvailableMachineTypes get all machine types that can be requested from the cloud provider.
func (asrProvider *autoScalrCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	return asrProvider.awsProvider.GetAvailableMachineTypes()
}

// NewNodeGroup builds a theoretical node group based on the node definition provided. The node group is not automatically
// created on the cloud provider side. The node group is not returned by NodeGroups() until it is created.
<<<<<<< HEAD
func (asrProvider *autoScalrCloudProvider) NewNodeGroup(machineType string, labels map[string]string, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	awsNg, err := asrProvider.awsProvider.NewNodeGroup(machineType, labels, extraResources)
	if err != nil {
		return awsNg, err
	} else {
		// wrap in asrNode
		return BuildAutoScalrNodeGroup(awsNg), err
	}
=======
func (asrProvider *autoScalrCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string,
	extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	return asrProvider.awsProvider.NewNodeGroup(machineType, labels, systemLabels, extraResources )
>>>>>>> cluster-autoscaler-release-1.2
}

// GetResourceLimiter returns struct containing limits (max, min) for resources (cores, memory etc.).
func (asrProvider *autoScalrCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	return asrProvider.awsProvider.GetResourceLimiter()
}

<<<<<<< HEAD
// Cleanup stops the go routine that is handling the current view of the ASGs in the form of a cache
func (asrProvider *autoScalrCloudProvider) Cleanup() error {
	return asrProvider.awsProvider.Cleanup()
}
=======
var launchTime = time.Now()
>>>>>>> cluster-autoscaler-release-1.2

// Refresh is called before every main loop and can be used to dynamically update cloud provider state.
// In particular the list of node groups returned by NodeGroups can change as a result of CloudProvider.Refresh().
func (asrProvider *autoScalrCloudProvider) Refresh() error {
<<<<<<< HEAD
	return asrProvider.awsProvider.Refresh()
=======
	err := asrProvider.awsProvider.Refresh()
	if err != nil {
		glog.Errorf("Failed to refresh cloud provider config: %v", err)
		return err
	}
	var execTime = time.Now()
	var elapsedTime = execTime.Sub(launchTime)
	glog.V(4).Info("Running for ", elapsedTime.Hours(), " hours" )
	if elapsedTime.Hours() > 24 {
		glog.V(4).Info("Running over 24 hours, exiting to force restart.")
		os.Exit(0)
	}

	depFlag := os.Getenv("ANAYLZE_DEPLOYMENTS")
	if depFlag != "false" {
		err = asrProvider.CollectClusterState()
	}
	return err
}

func (asrProvider *autoScalrCloudProvider) CollectClusterState() error {
	glog.V(4).Info("Starting CollectClusterState.")
	kubeClient := asrProvider.createKubeClient()
	nodeList, err := kubeClient.CoreV1().Nodes().List(metav1.ListOptions{})

	//allNodes, err := allNodeLister.List()
	if err != nil {
		glog.Errorf("Failed to list all nodes: %v", err)
	}
	// print out node names
	for _, aNode := range nodeList.Items {
		glog.V(4).Info("Node: ", aNode.Name, "Alloc mem: ", aNode.Status.Allocatable.Memory())
	}

	//podList, err := kubeClient.CoreV1().Pods(apiv1.NamespaceAll).List(metav1.ListOptions{})
	//if err != nil {
	//	glog.Errorf("Failed to list all pods: %v", err)
	//}
	// print out pods names
	//for _, aPod := range podList.Items {
	//	glog.V(4).Info("Pod: ", aPod.Name)
	//}
	// not needed now

	servList, err := kubeClient.AppsV1().Deployments(apiv1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("Failed to list all services: %v", err)
	}
	// print out service names
	for _, aServ := range servList.Items {
		glog.V(4).Info("Service: ", aServ.Name, " replicas: ", aServ.Status.Replicas)
	}
	state := &AutoScalrClusterState{
		AsrToken:    os.Getenv("AUTOSCALR_API_KEY"),
		AwsRegion:    os.Getenv("AWS_REGION"),
		AutoScalingGroupName:    os.Getenv("AUTOSCALING_GROUP_NAME"),
		Deployments: servList.Items,
		Nodes: nodeList.Items,
	}
	rc, err := SendClusterState(state, kubeClient)
	glog.V(4).Info("AutoScalrClusterState returned: ", rc)
	if err != nil {
		glog.Errorf("Error in SendClusterState: %v", err)
	} else {
		err = fmt.Errorf("CollectClusterState.Completed")
	}
	return err
}

func (asrProvider *autoScalrCloudProvider) createKubeClient() kube_client.Interface {
	url, err := url.Parse("")
	kubeConfig, err := config.GetKubeClientConfig(url)
	if err != nil {
		glog.Fatalf("Failed to build Kubernetes client configuration: %v", err)
	}
	return kube_client.NewForConfigOrDie(kubeConfig)
>>>>>>> cluster-autoscaler-release-1.2
}

// asrNodeGroup implements NodeGroup interface, defaulting to pass through to awsNodeGroup object
type asrNodeGroup struct {
	awsNodeGroup			cloudprovider.NodeGroup
}

func BuildAutoScalrNodeGroup(aNode cloudprovider.NodeGroup) (cloudprovider.NodeGroup) {
	asrNG := &asrNodeGroup{
		awsNodeGroup: aNode,
	}
	return asrNG
}

func (asrNG *asrNodeGroup) MaxSize() int {
	//glog.V(0).Infof("AsrNodeGroup::MaxSize")
	return asrNG.awsNodeGroup.MaxSize()
}

func (asrNG *asrNodeGroup) MinSize() int {
	//glog.V(0).Infof("AsrNodeGroup::MinSize")
	return asrNG.awsNodeGroup.MinSize()
}

func (asrNG *asrNodeGroup) TargetSize() (int, error) {
	//glog.V(0).Infof("AsrNodeGroup::TargetSize")
	app, err := appDefRead()
	tSize := 0
	if err != nil {
		glog.V(0).Infof("Received error from appDefRead: %s", err.Error())
	}
	if app != nil {
		numVcpu := numVCpusBaseType()
		tSize = app.TargetCapacity / numVcpu
	}
	glog.V(0).Infof("Returning TargetSize: %d",tSize)
	return tSize, err
}

func (asrNG *asrNodeGroup) IncreaseSize(delta int) error {
	glog.V(0).Infof("AsrNodeGroup::IncreaseSize delta: %v",delta)
	currSize, err := asrNG.TargetSize()
	if err != nil {
		glog.V(0).Infof("TargetSize returned error: %v", err.Error())
	} else {
		//glog.V(0).Infof("currSize: %v",currSize)
		numVcpu := numVCpusBaseType()
		//glog.V(0).Infof("numVcpu: %v",numVcpu)
		newTarget := (currSize + delta) * numVcpu
		glog.V(0).Infof("new vCpu target: %v",newTarget)
		err = appDefUpdate(newTarget)
	}
	return err
}

func (asrNG *asrNodeGroup) DeleteNodes(nodes []*apiv1.Node) error {
	//glog.V(0).Infof("AsrNodeGroup::DeleteNodes")
	numNodesToDelete := len(nodes)
	glog.V(0).Infof("Deleting %v nodes",numNodesToDelete)

	nodeIds := make([]string, 0, len(nodes))
	for _, node := range nodes {
		provId := node.Spec.ProviderID
		instId := InstanceIdFromProviderId(provId)
		glog.V(0).Infof("Deleting instance id: %v",instId)
		//glog.V(0).Infof("node.Spec: %v",node.Spec.String())
		nodeIds = append(nodeIds, instId)
	}
	err := appDefDeleteNodes(0, nodeIds)
	if err != nil {
		glog.V(0).Infof("Received error from appDefDeleteNodes: %s", err.Error())
	}
	return err
}

func (asrNG *asrNodeGroup) DecreaseTargetSize(delta int) error {
	//glog.V(0).Infof("AsrNodeGroup::DecreaseTargetSize")
	return nil
	//return cloudprovider.ErrNotImplemented
	//return asrNG.awsNodeGroup.DecreaseTargetSize(delta)
}

func (asrNG *asrNodeGroup) Id() string {
	//glog.V(0).Infof("AsrNodeGroup::Id")
	return asrNG.awsNodeGroup.Id()
}

func (asrNG *asrNodeGroup) Debug() string {
	return asrNG.awsNodeGroup.Debug()
}

func (asrNG *asrNodeGroup) Nodes() ([]string, error) {
<<<<<<< HEAD
	//glog.V(0).Infof("AsrNodeGroup::Nodes")
=======
	glog.V(0).Infof("AsrNodeGroup::Nodes")
>>>>>>> cluster-autoscaler-release-1.2
	return asrNG.awsNodeGroup.Nodes()
}

func (asrNG *asrNodeGroup) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
	//glog.V(0).Infof("AsrNodeGroup::TemplateNodeInfo")
	return asrNG.awsNodeGroup.TemplateNodeInfo()
}

func (asrNG *asrNodeGroup) Exist() bool {
	//glog.V(0).Infof("AsrNodeGroup::Exist")
	app, err := appDefRead()
	exists := false
	if err != nil {
		glog.V(0).Infof("Received error from appDefRead: %s", err.Error())
	}
	if app != nil {
		exists = true
	}
	//glog.V(0).Infof("Returning exists: %v",exists)
	return exists
	//return asrNG.awsNodeGroup.Exist()
}

func (asrNG *asrNodeGroup) Create() error {
	glog.V(0).Infof("AsrNodeGroup::Create")
	return appDefCreate()
	//return asrNG.awsNodeGroup.Create()
}

func (asrNG *asrNodeGroup) Delete() error {
	glog.V(0).Infof("AsrNodeGroup::Delete")
	return appDefDelete()
	//return asrNG.awsNodeGroup.Delete()
}

func (asrNG *asrNodeGroup) Autoprovisioned() bool {
	//glog.V(0).Infof("AsrNodeGroup::Autoprovisioned")
	return false
}
