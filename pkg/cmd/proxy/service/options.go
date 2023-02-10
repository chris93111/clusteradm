// Copyright Contributors to the Open Cluster Management project
package service

import (
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"
	//"sigs.k8s.io/kustomize/kyaml/errors"
)

// Options: only support use in-cluster certificates
type Options struct {
	//ClusteradmFlags: The generic options from the clusteradm cli-runtime.
	ClusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags

	cluster               string
	namespace             string
	service               string
	port                  int32
	secure                bool
	managedServiceAccount string
	kubectlArgs           string
}

func newOptions(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags) *Options {
	return &Options{
		ClusteradmFlags: clusteradmFlags,
	}
}

func (o *Options) validate() error {
	return nil
}
