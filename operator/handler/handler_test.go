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
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tailingsidecarv1 "github.com/SumoLogic/tailing-sidecar/operator/api/v1"
	admv1 "k8s.io/api/admission/v1beta1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestPodExtender(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PodExtender Suite")
}

var _ = Describe("handler", func() {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	Context("PodExtender.Handle", func() {
		testEnv := &envtest.Environment{
			CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
		}
		cfg, err := testEnv.Start()
		It("", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg).ToNot(BeNil())
		})

		err = tailingsidecarv1.AddToScheme(scheme.Scheme)
		It("", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		k8sClient, clientErr := client.New(cfg, client.Options{Scheme: scheme.Scheme})
		It("", func() {
			Expect(clientErr).ToNot(HaveOccurred())
		})

		decoder, err := admission.NewDecoder(scheme.Scheme)
		It("creates decoder without any errors", func() {
			Expect(err).To(BeNil())

		})

		podExtender := PodExtender{
			Client:              k8sClient,
			TailingSidecarImage: "tailing-sidecar-image:test",
			decoder:             decoder,
		}

		When("request does not contain any object", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(``),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("rejects request as decoder returns an error", func() {
				Expect(resp.Allowed).To(BeFalse())
				Expect(resp.Result.Code).Should(Equal(int32(http.StatusBadRequest)))
				s, _ := json.MarshalIndent(resp.Patches, "", "\t")
				fmt.Printf("patch: %+v\n", string(s))
			})
		})

		When("request contains empty json", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{}`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)

			It("returns empty patch as extendPod function didn't find tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				s, _ := json.MarshalIndent(resp.Patches, "", "\t")
				fmt.Printf("patch: %+v\n", string(s))
			})
		})

		When("Pod with null metadata is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": null,
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox"
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns empty patch as extendPod function didn't find tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				// s, _ := json.MarshalIndent(resp.Patches, "", "\t")
				// fmt.Printf("patch: %+v\n", string(s))
			})
		})

		When("Pod with empty metadata is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {},
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
			It("returns empty patch as extendPod function didn't find tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				// Expect(resp.Patches).To(BeEmpty())
				// s, _ := json.MarshalIndent(resp.Patches, "", "\t")
				// fmt.Printf("patch: %+v\n", string(s))
			})
		})

		When("Pod with empty annotation is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {
										"name": "simple",
									  	"annotations": {}
									},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox"
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns empty patch as extendPod function didn't find tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				// Expect(resp.Patches).To(BeEmpty())
				// s, _ := json.MarshalIndent(resp.Patches, "", "\t")
				// fmt.Printf("patch: %+v\n", string(s))
			})
		})

		When("Pod with null annotation is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
									"apiVersion": "v1",
									"kind": "Pod",
									"metadata": {
										"name": "simple",
									  	"annotations": null
									},
									"spec": {
									  "containers": [
										{
										  "name": "busybox",
										  "image": "busybox"
										}
									  ]
									}
								  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns empty patch as extendPod function didn't find tailing-sidecar annotation", func() {
				Expect(resp.Allowed).To(BeTrue())
				// s, _ := json.MarshalIndent(resp.Patches, "", "\t")
				// fmt.Printf("patch: %+v\n", string(s))
			})
		})

		When("Pod with correct configuration is created", func() {
			request := admission.Request{
				AdmissionRequest: admv1.AdmissionRequest{
					Operation: admv1.Create,
					Object: runtime.RawExtension{
						Raw: []byte(`{
							"apiVersion": "v1",
							"kind": "Pod",
							"metadata": {
							  "name": "pod-with-annotations",
							  "namespace": "tailing-sidecar-system",
							  "annotations": {
								"tailing-sidecar": "varlog:/var/log/example0.log;varlog:/var/log/example1.log"
							  }
							},
							"spec": {
							  "containers": [
								{
								  "name": "count",
								  "image": "busybox",
								  "args": [
									"/bin/sh",
									"-c",
									"i=0; while true; do\n  echo \"example0: $i $(date)\" >> /var/log/example0.log;\n  echo \"example1: $i $(date)\" >> /var/log/example1.log;\n  echo \"example2: $i $(date)\" >> /varconfig/log/example2.log;\n  i=$((i+1));\n  sleep 1;\ndone\n"
								  ],
								  "volumeMounts": [
									{
									  "name": "varlog",
									  "mountPath": "/var/log"
									},
									{
									  "name": "varlogconfig",
									  "mountPath": "/varconfig/log"
									}
								  ]
								}
							  ],
							  "volumes": [
								{
								  "name": "varlog",
								  "emptyDir": {}
								},
								{
								  "name": "varlogconfig",
								  "emptyDir": {}
								}
							  ]
							}
						  }`),
					},
				},
			}

			resp := podExtender.Handle(context.Background(), request)
			It("returns patch with tailing sidecar containers", func() {
				Expect(resp.Allowed).To(BeTrue())
				Expect(resp.Patches).NotTo(BeEmpty())
				//nodesJSON, err := json.Marshal(nodes)
				//assert.Nil(t, err)
				//assert.JSONEq(t, expectedNodes, string(nodesJSON))
				// s, _ := json.MarshalIndent(resp, "", "\t")
				// fmt.Printf("patch: %+v\n", string(s))
				//fmt.Println("HTTP code:", resp.Result.Code)
			})
		})
	})
})
