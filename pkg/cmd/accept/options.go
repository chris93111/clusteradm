// Copyright Contributors to the Open Cluster Management project
package accept

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"
)

type Options struct {
	//ClusteradmFlags: The generic options from the clusteradm cli-runtime.
	ClusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags
	//A list of comma separated cluster names
	Clusters string
	//Wait to wait for managedcluster and CSR
	Wait bool
	//If true the csr will approve directly and check of requester will skip.
	SkipApproveCheck bool

	Values Values

	Streams genericclioptions.IOStreams
}

//Values: The values used in the template
type Values struct {
	Clusters []string
}

func NewOptions(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		ClusteradmFlags: clusteradmFlags,
		Streams:         streams,
	}
}
