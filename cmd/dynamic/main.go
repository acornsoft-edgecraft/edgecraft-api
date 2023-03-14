package main

import (
	"fmt"
	"log"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// var (
// 	runtimeClassGVR = schema.GroupVersionResource{
// 		Group:    "node.k8s.io",
// 		Version:  "v1alpha1",
// 		Resource: "runtimeclasses",
// 	}
// )

func main() {
	log.Print("Loading client config")
	config, err := clientcmd.BuildConfigFromFlags("", userConfig())
	errExit("Failed to load client conifg", err)

	log.Print("Loading dynamic client")
	client, err := dynamic.NewForConfig(config)
	errExit("Failed to create client", err)

	//RegisterRuntimeClassCRD(config)
	//CreateSampleRuntimeClasses(client)
	//PrintRuntimeHandlers(client)
	PrintResources(client, "os-ka-upd")

	PrintStatus(client, "default", "os-ka-upd-m", "controlplane.cluster.x-k8s.io", "v1alpha3", "kubeadmcontrolplanes")
}

// func RegisterRuntimeClassCRD(config *rest.Config) {
// 	apixClient, err := apixv1beta1client.NewForConfig(config)
// 	errExit("Failed to load apiextensions client", err)

// 	crds := apixClient.CustomResourceDefinitions()

// 	const (
// 		dns1123LabelFmt        string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
// 		dns1123SubdomainFmt    string = dns1123LabelFmt + "(\\." + dns1123LabelFmt + ")*"
// 		dns1123SubdomainRegexp string = "^" + dns1123SubdomainFmt + "$"
// 	)
// 	runtimeClassCRD := &apixv1beta1.CustomResourceDefinition{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: "runtimeclasses.node.k8s.io",
// 		},
// 		Spec: apixv1beta1.CustomResourceDefinitionSpec{
// 			Group:   "node.k8s.io",
// 			Version: "v1alpha1",
// 			Versions: []apixv1beta1.CustomResourceDefinitionVersion{{
// 				Name:    "v1alpha1",
// 				Served:  true,
// 				Storage: true,
// 			}},
// 			Names: apixv1beta1.CustomResourceDefinitionNames{
// 				Plural:   "runtimeclasses",
// 				Singular: "runtimeclass",
// 				Kind:     "RuntimeClass",
// 			},
// 			Scope: apixv1beta1.ClusterScoped,
// 			Validation: &apixv1beta1.CustomResourceValidation{
// 				OpenAPIV3Schema: &apixv1beta1.JSONSchemaProps{
// 					Properties: map[string]apixv1beta1.JSONSchemaProps{
// 						"spec": {
// 							Properties: map[string]apixv1beta1.JSONSchemaProps{
// 								"runtimeHandler": {
// 									Type:    "string",
// 									Pattern: dns1123SubdomainRegexp,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	log.Print("Registering RuntimeClass CRD")
// 	_, err = crds.Create(runtimeClassCRD)
// 	if err != nil {
// 		if errors.IsAlreadyExists(err) {
// 			log.Print("RuntimeClass CRD already registered")
// 		} else {
// 			errExit("Failed to create RuntimeClass CRD", err)
// 		}
// 	}
// }

// func CreateSampleRuntimeClasses(client dynamic.Interface) {
// 	res := client.Resource(runtimeClassGVR)

// 	rcs := map[string]string{
// 		"native":  "runc",
// 		"sandbox": "gvisor",
// 		"vm":      "kata-containers",
// 		"foo":     "bar",
// 	}
// 	for name, handler := range rcs {
// 		log.Printf("Creating RuntimeClass %s", name)
// 		rc := NewRuntimeClass(name, handler)
// 		_, err := res.Create(rc, metav1.CreateOptions{})
// 		errExit(fmt.Sprintf("Failed to create RuntimeClass %#v", rc), err)
// 	}
// }

// func NewRuntimeClass(name, handler string) *unstructured.Unstructured {
// 	return &unstructured.Unstructured{
// 		Object: map[string]interface{}{
// 			"kind":       "RuntimeClass",
// 			"apiVersion": runtimeClassGVR.Group + "/v1alpha1",
// 			"metadata": map[string]interface{}{
// 				"name": name,
// 			},
// 			"spec": map[string]interface{}{
// 				"runtimeHandler": handler,
// 			},
// 		},
// 	}
// }

// func PrintRuntimeHandlers(client dynamic.Interface) {
// 	PrintResourceField(client, runtimeClassGVR, "spec", "runtimeHandler")
// }

// func PrintResourceField(client dynamic.Interface, gvr schema.GroupVersionResource, fldPath ...string) {
// 	rs := fmt.Sprintf("%s/%s", gvr.Group, gvr.Resource)
// 	log.Printf("Listing %s objects", rs)
// 	res := client.Resource(gvr)
// 	list, err := res.List(metav1.ListOptions{})
// 	errExit("Failed to list "+rs+" objects", err)

// 	log.Printf("Printing %s.%s", rs, strings.Join(fldPath, "."))
// 	output := make(map[string]string)
// 	for _, item := range list.Items {
// 		name := item.GetName()
// 		fld, exists, err := unstructured.NestedString(item.Object, fldPath...)
// 		if err != nil {
// 			log.Printf("Error reading %s for %s: %v", strings.Join(fldPath, "."), name, err)
// 			continue
// 		}
// 		if !exists {
// 			fld = "[NOT FOUND]"
// 		}
// 		output[name] = fld
// 	}

// 	for name, fld := range output {
// 		fmt.Printf("  %-10s  -->  %-10s\n", name, fld)
// 	}
// }

func PrintResources(client dynamic.Interface, clusterName string) {
	gvr := schema.GroupVersionResource{Group: "infrastructure.cluster.x-k8s.io", Version: "v1alpha3", Resource: "openstackclusters"}
	rs := fmt.Sprintf("%s/%s", gvr.Group, gvr.Resource)
	log.Printf("Listing %s objects", rs)
	res := client.Resource(gvr).Namespace("default")
	// log.Printf("Printing %s.%s", rs, "status.ready")

	// 단일 정보
	item, err := res.Get(context.TODO(), clusterName, metav1.GetOptions{})
	errExit("Failed to list "+rs+" objects", err)
	fld, exists, err := unstructured.NestedBool(item.Object, "status", "ready")
	if err != nil {
		log.Printf("Error reading %s for %s: %v", "status.ready", item.GetName(), err)
	}
	if !exists {
		fld = false
	}
	fmt.Printf("  %-10s  -->  %-10v\n", item.GetName(), fld)
}

// PrintStatus - Print Custom Resource Obejct's Status Info
func PrintStatus(client dynamic.Interface, namespace, resourceName, group, version, resource string) {
	gvr := schema.GroupVersionResource{Group: group, Version: version, Resource: resource}
	rs := fmt.Sprintf("%s/%s", gvr.Group, gvr.Resource)
	log.Printf("Listing %s objects", rs)

	res := client.Resource(gvr).Namespace(namespace)
	item, err := res.Get(context.TODO(), resourceName, metav1.GetOptions{})
	errExit("Failed to list "+rs+" objects", err)

	flds, found, err := unstructured.NestedSlice(item.Object, "status", "conditions")
	if err != nil {
		log.Printf("Error reading %s for %s: %v", "status.conditions", item.GetName(), err)
	} else if !found {
		log.Printf("Not found %s for %s", "status.conditions", item.GetName())
	}

	for _, itemRaw := range flds {
		item := itemRaw.(map[string]interface{})
		value, found, err := unstructured.NestedString(item, "type")
		if !found || err != nil {
			log.Printf("Not found type or error: %v\n", err)
			continue
		}
		fmt.Printf("Found type: %s\n", value)

		status, found, err := unstructured.NestedString(item, "status")
		if !found || err != nil {
			log.Printf("Not found status or error: %v\n", err)
			continue
		}
		fmt.Printf("Found status: %s\n", status)

		if value == "MachinesSpecUpToDate" {
			fmt.Printf("Target :: %s - %v\n", value, status)
		}
	}
}

func errExit(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %#v", msg, err)
	}
}

func userConfig() string {
	usr, err := user.Current()
	errExit("Failed to get current user", err)

	return filepath.Join(usr.HomeDir, ".kube", "config")
}
