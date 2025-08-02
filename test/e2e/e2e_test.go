/*
Copyright 2024.

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

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/seaweedfs/seaweedfs-operator/test/utils"
)

const namespace = "seaweedfs-operator-system"
const deploymentName = "seaweedfs-operator-controller-manager"

var _ = Describe("controller", func() {

	Context("Operator", func() {
		It("should run successfully", func() {
			kubeconfig := config.GetConfigOrDie()
			clientset, err := utils.GetClientset(kubeconfig)

			deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
			if err != nil {
				panic(err.Error())
			}

			// Get the pods under the deployment
			selector := deployment.Spec.Selector.MatchLabels
			labelSelector := metav1.LabelSelector{MatchLabels: selector}

			pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: metav1.FormatLabelSelector(&labelSelector),
			})

			if err != nil {
				panic(err.Error())
			}

			stopCh := make(chan struct{}, 1)
			readyCh := make(chan struct{})
			localPort, err := utils.GetFreePort()
			localPortStr := strconv.Itoa(localPort)

			// Call the function to run port forward
			err = utils.RunPortForward(kubeconfig, namespace, pods.Items[0].Name, []string{fmt.Sprintf("%s:8081", localPortStr)}, stopCh, readyCh)
			if err != nil {
				panic(err.Error())
			}
			<-readyCh

			readyzURL := fmt.Sprintf("http://localhost:%s/readyz", localPortStr)
			client := resty.New().SetTimeout(5 * time.Second)
			resp, err := client.R().Get(readyzURL)
			if err != nil {
				panic(err.Error())
			}

			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp.Body()).To(Equal([]uint8{'o', 'k'}))

			close(stopCh)
			<-stopCh
		})
	})
})
