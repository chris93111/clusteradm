// Copyright Contributors to the Open Cluster Management project
package work

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	workclient "open-cluster-management.io/api/client/work/clientset/versioned"
	workapiv1 "open-cluster-management.io/api/work/v1"
	"open-cluster-management.io/clusteradm/pkg/helpers/printer"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 1 {
		return fmt.Errorf("can only specify one manifestwork")
	}

	if len(args) == 1 {
		o.workName = args[0]
	}

	o.printer.Competele()

	return nil
}

func (o *Options) validate() (err error) {
	err = o.ClusteradmFlags.ValidateHub()
	if err != nil {
		return err
	}

	if len(o.cluster) == 0 {
		return fmt.Errorf("cluster name must be specified")
	}

	err = o.printer.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (o *Options) run() (err error) {
	restConfig, err := o.ClusteradmFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}
	clusterClient, err := clusterclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	workClient, err := workclient.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	_, err = clusterClient.ClusterV1().ManagedClusters().Get(context.TODO(), o.cluster, metav1.GetOptions{})
	if err != nil {
		return err
	}

	var workList *workapiv1.ManifestWorkList
	if len(o.workName) == 0 {
		workList, err = workClient.WorkV1().ManifestWorks(o.cluster).List(context.TODO(), metav1.ListOptions{})
	} else {
		workList, err = workClient.WorkV1().ManifestWorks(o.cluster).List(context.TODO(), metav1.ListOptions{
			FieldSelector: fmt.Sprintf("name=%s", o.workName),
		})
	}
	if err != nil {
		return err
	}

	o.printer.WithTreeConverter(o.convertToTree).WithTableConverter(o.converToTable)

	return o.printer.Print(o.Streams, workList)
}

func (o *Options) convertToTree(obj runtime.Object, tree *printer.TreePrinter) *printer.TreePrinter {
	if workList, ok := obj.(*workapiv1.ManifestWorkList); ok {
		for _, work := range workList.Items {
			cluster, number, applied, available := getFileds(work)
			mp := make(map[string]interface{})
			mp[".Cluster"] = cluster
			mp[".Number of Manifests"] = number
			mp[".Applied"] = applied
			mp[".Available"] = available

			tree.AddFileds(work.Name, &mp)
		}
	}
	return tree
}

func (o *Options) converToTable(obj runtime.Object) *metav1.Table {
	table := &metav1.Table{
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Cluster", Type: "string"},
			{Name: "Number Of Manifests", Type: "integer"},
			{Name: "Applied", Type: "string"},
			{Name: "Available", Type: "string"},
		},
		Rows: []metav1.TableRow{},
	}

	if workList, ok := obj.(*workapiv1.ManifestWorkList); ok {
		for _, work := range workList.Items {
			cluster, number, applied, available := getFileds(work)
			row := metav1.TableRow{
				Cells:  []interface{}{work.Name, cluster, number, applied, available},
				Object: runtime.RawExtension{Object: &work},
			}

			table.Rows = append(table.Rows, row)
		}
	}

	return table
}

func getFileds(work workapiv1.ManifestWork) (cluster string, number int, applied, available string) {
	cluster = work.Namespace
	number = len(work.Spec.Workload.Manifests)

	appliedCond := meta.FindStatusCondition(work.Status.Conditions, workapiv1.WorkApplied)
	if appliedCond != nil {
		applied = string(appliedCond.Status)
	}

	availableCond := meta.FindStatusCondition(work.Status.Conditions, workapiv1.WorkAvailable)
	if availableCond != nil {
		available = string(availableCond.Status)
	}

	return
}
