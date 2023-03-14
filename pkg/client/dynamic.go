/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// DynamicClient - Dynamic client type
type DynamicClient struct {
	config       *rest.Config                // Rest config
	resource     schema.GroupVersionResource // resource gvr definition
	namespace    string                      // namespace
	hasNamespace bool                        // namespace 설정 여부
}

// SetNamespace - Namespace 설정
func (dc *DynamicClient) SetNamespace(ns string) {
	dc.namespace = ns
	dc.hasNamespace = (ns != "")
}

// List - ListOptions에 따른 List 처리
func (dc *DynamicClient) List(opts v1.ListOptions) (res *unstructured.UnstructuredList, err error) {
	// Get dynamic interface
	di, err := dynamic.NewForConfig(dc.config)
	if err != nil {
		return
	}

	if dc.hasNamespace {
		res, err = di.Resource(dc.resource).Namespace(dc.namespace).List(context.TODO(), opts)
	} else {
		res, err = di.Resource(dc.resource).List(context.TODO(), opts)
	}

	return
}

// Get - 이름과 옵션에 따른 Get 처리
func (dc *DynamicClient) Get(name string, opts v1.GetOptions) (res *unstructured.Unstructured, err error) {
	// Get dynamic interface
	di, err := dynamic.NewForConfig(dc.config)
	if err != nil {
		return
	}

	if dc.hasNamespace {
		res, err = di.Resource(dc.resource).Namespace(dc.namespace).Get(context.TODO(), name, opts)
	} else {
		res, err = di.Resource(dc.resource).Get(context.TODO(), name, opts)
	}

	return
}

// Watch - 옵션에 따른 Watch 처리
func (dc *DynamicClient) Watch(opts v1.ListOptions) (wi watch.Interface, err error) {
	// Get dynamic interface
	di, err := dynamic.NewForConfig(dc.config)
	if err != nil {
		return
	}
	opts.Watch = true

	if dc.hasNamespace {
		wi, err = di.Resource(dc.resource).Namespace(dc.namespace).Watch(context.TODO(), opts)
	} else {
		wi, err = di.Resource(dc.resource).Watch(context.TODO(), opts)
	}

	return
}

// Delete - 이름과 옵션에 따른 Delete 처리
func (dc *DynamicClient) Delete(name string, opts v1.DeleteOptions) (err error) {
	// Get dynamic interface
	di, err := dynamic.NewForConfig(dc.config)
	if err != nil {
		return
	}

	if dc.hasNamespace {
		err = di.Resource(dc.resource).Namespace(dc.namespace).Delete(context.TODO(), name, opts)
	} else {
		err = di.Resource(dc.resource).Delete(context.TODO(), name, opts)
	}

	return
}

// MultiplePost - 데이터 스트림과 갱신여부에 따른 POST 처리 (Multiple Resource)
func (dc *DynamicClient) MultiplePost(payload io.Reader) (resources []*unstructured.Unstructured, err error) {
	// 데이터 스트림 Decode
	decoder := yaml.NewYAMLOrJSONDecoder(payload, 4096)

	for {
		// decode payload
		data := &unstructured.Unstructured{}
		err = decoder.Decode(data)
		if err == io.EOF || err != nil {
			// No content or err
			break
		}
		if data.Object == nil {
			continue
		}

		// get version, kind
		version := data.GetAPIVersion()
		kind := data.GetKind()

		gv, err := schema.ParseGroupVersion(version)
		if err != nil {
			gv = schema.GroupVersion{Version: version}
		}

		// 설정에 따르는 Discovery 클라이언트 생성
		discoveryClient, err := discovery.NewDiscoveryClientForConfig(dc.config)
		if err != nil {
			return resources, err
		}

		// Version에 맞는 리소스 리스트 조회
		apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(version)
		if err != nil {
			return resources, err
		}

		apiResources := apiResourceList.APIResources

		// Kind에 해당하는 리소스 검색
		var resource *v1.APIResource
		for _, apiResource := range apiResources {
			if apiResource.Kind == kind && !strings.Contains(apiResource.Name, "/") {
				resource = &apiResource
				break
			}
		}

		if resource == nil {
			err = fmt.Errorf("unknown resource kind: %s", kind)
			return resources, err
		}

		// Get dynamic interface
		di, err := dynamic.NewForConfig(dc.config)
		if err != nil {
			return resources, err
		}

		dc.resource = schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource.Name}
		dc.namespace = data.GetNamespace()

		// Check resource
		var res *unstructured.Unstructured
		_, err = di.Resource(dc.resource).Namespace(dc.namespace).Get(context.TODO(), data.GetName(), v1.GetOptions{})
		if err != nil {
			if utils.CheckK8sNotFound(err) {
				// Create
				if resource.Namespaced {
					res, err = di.Resource(dc.resource).Namespace(dc.namespace).Create(context.TODO(), data, v1.CreateOptions{})
					if err != nil {
						return resources, err
					}
				} else {
					res, err = di.Resource(dc.resource).Create(context.TODO(), data, v1.CreateOptions{})
					if err != nil {
						return resources, err
					}
				}
			} else {
				return resources, nil
			}
		} else {
			// Update
			jsonData, err := json.Marshal(data)
			if err != nil {
				return resources, err
			}
			res, err = di.Resource(dc.resource).Namespace(dc.namespace).Patch(context.TODO(), data.GetName(), types.MergePatchType, jsonData, v1.PatchOptions{FieldManager: "server-side-apply"})
			if err != nil {
				return resources, err
			}
		}

		resources = append(resources, res)
	}

	return resources, nil
}

// Post - 데이터 스트림과 갱신여부에 따른 POST 처리 (Single Resource)
func (dc *DynamicClient) Post(payload io.Reader, isUpdate bool) (res *unstructured.Unstructured, err error) {
	// 데이터 스트림 Decode
	decoder := yaml.NewYAMLOrJSONDecoder(payload, 4096)

	for {
		// decode payload
		data := &unstructured.Unstructured{}
		if err = decoder.Decode(data); err != nil {
			return
		}

		// get version, kind
		version := data.GetAPIVersion()
		kind := data.GetKind()

		gv, err := schema.ParseGroupVersion(version)
		if err != nil {
			gv = schema.GroupVersion{Version: version}
		}

		// 설정에 따르는 Discovery 클라이언트 생성
		discoveryClient, err := discovery.NewDiscoveryClientForConfig(dc.config)
		if err != nil {
			return res, err
		}

		// Version에 맞는 리소스 리스트 조회
		apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(version)
		if err != nil {
			return res, err
		}

		apiResources := apiResourceList.APIResources

		// Kind에 해당하는 리소스 검색
		var resource *v1.APIResource
		for _, apiResource := range apiResources {
			if apiResource.Kind == kind && !strings.Contains(apiResource.Name, "/") {
				resource = &apiResource
				break
			}
		}

		if resource == nil {
			err = fmt.Errorf("unknown resource kind: %s", kind)
			return res, err
		}

		// Get dynamic interface
		di, err := dynamic.NewForConfig(dc.config)
		if err != nil {
			return res, err
		}

		dc.resource = schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource.Name}
		dc.namespace = data.GetNamespace()

		// Update인 경우는 resourceVersion을 조회하고 수정
		if isUpdate {
			targetRes, err := di.Resource(dc.resource).Namespace(dc.namespace).Get(context.TODO(), data.GetName(), v1.GetOptions{})
			if err != nil {
				return res, err
			}
			data.SetResourceVersion(targetRes.GetResourceVersion())
			if resource.Namespaced {
				//lint:ignore SA4006 Ignore unused variable
				res, err = di.Resource(dc.resource).Namespace(dc.namespace).Update(context.TODO(), data, v1.UpdateOptions{})
			} else {
				//lint:ignore SA4006 Ignore
				res, err = di.Resource(dc.resource).Update(context.TODO(), data, v1.UpdateOptions{})
			}
		} else {
			if resource.Namespaced {
				res, err = di.Resource(dc.resource).Namespace(dc.namespace).Create(context.TODO(), data, v1.CreateOptions{})
			} else {
				res, err = di.Resource(dc.resource).Create(context.TODO(), data, v1.CreateOptions{})
			}
		}

		//lint:ignore SA4004 Ignore unconditional terminated loop
		return res, err
	}
}

// Patch - 지정한 이름과 형식 및 데이터에 따라 Patch 처리
func (dc *DynamicClient) Patch(name string, patchType types.PatchType, payload io.Reader, opts v1.PatchOptions) (res *unstructured.Unstructured, err error) {
	data, err := io.ReadAll(payload)
	if err != nil {
		return res, err
	}

	// Get dynamic interface
	di, err := dynamic.NewForConfig(dc.config)
	if err != nil {
		return res, err
	}

	if dc.hasNamespace {
		res, err = di.Resource(dc.resource).Namespace(dc.namespace).Patch(context.TODO(), name, patchType, data, opts)
	} else {
		res, err = di.Resource(dc.resource).Patch(context.TODO(), name, patchType, data, opts)
	}

	return
}

// NewDynamicClient - Config 기준의 Restfull Client 생성
func NewDynamicClient(config *rest.Config) *DynamicClient {
	return &DynamicClient{
		config:       config,
		hasNamespace: false,
	}
}

// NewDynamicClientSchema - Config, Group, Version, Resource 기준의 Restful Client 생성
// ex. schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
func NewDynamicClientSchema(config *rest.Config, group, version, resource string) *DynamicClient {
	return &DynamicClient{
		config:       config,
		resource:     schema.GroupVersionResource{Group: group, Version: version, Resource: resource},
		hasNamespace: false,
	}
}
