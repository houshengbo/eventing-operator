/*
Copyright 2020 The Knative Authors

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "knative.dev/eventing-operator/pkg/apis/eventing/v1alpha1"
)

// FakeKnativeEventings implements KnativeEventingInterface
type FakeKnativeEventings struct {
	Fake *FakeOperatorV1alpha1
	ns   string
}

var knativeeventingsResource = schema.GroupVersionResource{Group: "operator.knative.dev", Version: "v1alpha1", Resource: "knativeeventings"}

var knativeeventingsKind = schema.GroupVersionKind{Group: "operator.knative.dev", Version: "v1alpha1", Kind: "KnativeEventing"}

// Get takes name of the knativeEventing, and returns the corresponding knativeEventing object, and an error if there is any.
func (c *FakeKnativeEventings) Get(name string, options v1.GetOptions) (result *v1alpha1.KnativeEventing, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(knativeeventingsResource, c.ns, name), &v1alpha1.KnativeEventing{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.KnativeEventing), err
}

// List takes label and field selectors, and returns the list of KnativeEventings that match those selectors.
func (c *FakeKnativeEventings) List(opts v1.ListOptions) (result *v1alpha1.KnativeEventingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(knativeeventingsResource, knativeeventingsKind, c.ns, opts), &v1alpha1.KnativeEventingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.KnativeEventingList{ListMeta: obj.(*v1alpha1.KnativeEventingList).ListMeta}
	for _, item := range obj.(*v1alpha1.KnativeEventingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested knativeEventings.
func (c *FakeKnativeEventings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(knativeeventingsResource, c.ns, opts))

}

// Create takes the representation of a knativeEventing and creates it.  Returns the server's representation of the knativeEventing, and an error, if there is any.
func (c *FakeKnativeEventings) Create(knativeEventing *v1alpha1.KnativeEventing) (result *v1alpha1.KnativeEventing, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(knativeeventingsResource, c.ns, knativeEventing), &v1alpha1.KnativeEventing{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.KnativeEventing), err
}

// Update takes the representation of a knativeEventing and updates it. Returns the server's representation of the knativeEventing, and an error, if there is any.
func (c *FakeKnativeEventings) Update(knativeEventing *v1alpha1.KnativeEventing) (result *v1alpha1.KnativeEventing, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(knativeeventingsResource, c.ns, knativeEventing), &v1alpha1.KnativeEventing{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.KnativeEventing), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeKnativeEventings) UpdateStatus(knativeEventing *v1alpha1.KnativeEventing) (*v1alpha1.KnativeEventing, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(knativeeventingsResource, "status", c.ns, knativeEventing), &v1alpha1.KnativeEventing{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.KnativeEventing), err
}

// Delete takes name of the knativeEventing and deletes it. Returns an error if one occurs.
func (c *FakeKnativeEventings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(knativeeventingsResource, c.ns, name), &v1alpha1.KnativeEventing{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeKnativeEventings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(knativeeventingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.KnativeEventingList{})
	return err
}

// Patch applies the patch and returns the patched knativeEventing.
func (c *FakeKnativeEventings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.KnativeEventing, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(knativeeventingsResource, c.ns, name, pt, data, subresources...), &v1alpha1.KnativeEventing{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.KnativeEventing), err
}
