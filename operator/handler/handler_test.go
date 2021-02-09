/*
Copyright 2021.

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

package handler

import (
	"context"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	admv1 "k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//https://github.com/influxdata/telegraf-operator/blob/e750174f2fbef5b04d3389c15ca66dd88fa4dc65/handler_test.go#L862

func TestPodExtender(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "PodExtender Suite")
}

var _ = Describe("PodExtender", func() {

	Context("PodExtender Handler", func() {

		dep := &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-deployment",
				Namespace: "ns1",
			},
		}

		decoder, err := admission.NewDecoder(runtime.NewScheme())
		It("decoder be created without errors", func() {
			Expect(err).To(BeNil())

		})

		podExtender := PodExtender{
			Client:              testclient.NewFakeClient(dep),
			TailingSidecarImage: "tailing-sidecar-image:test",
			decoder:             decoder,
		}

		When("When Pod with annotation is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Update,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {
										"name": "simple",
									  "annotations": {
										"telegraf.influxdata.com/port": "8080",
										"telegraf.influxdata.com/path": "/v1/metrics",
										"telegraf.influxdata.com/interval": "5s"
									  }
									},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox",
										  "args": [
											"sleep",
											"1000000"
										  ]
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns response with Pod extended by tailing sidecar containers", func() {
				Expect(resp).ToNot(BeNil())
				fmt.Println(resp.Allowed)
				fmt.Println("***")
				fmt.Println(resp.Patch)
				// assert.JSONEq(t, expectedNodes, string(nodesJSON))
			})
		})
	})
})
