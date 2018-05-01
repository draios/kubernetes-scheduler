/*
Copyright 2018 Sysdig.

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

package kubernetes

import "time"

type KubeReplicaSet struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata struct {
		Name              string            `json:"name"`
		Namespace         string            `json:"namespace"`
		SelfLink          string            `json:"selfLink"`
		UID               string            `json:"uid"`
		ResourceVersion   string            `json:"resourceVersion"`
		Generation        int               `json:"generation"`
		CreationTimestamp time.Time         `json:"creationTimestamp"`
		Labels            map[string]string `json:"labels"`
		Annotations       map[string]string `json:"annotations"`
		OwnerReferences []struct {
			APIVersion         string `json:"apiVersion"`
			Kind               string `json:"kind"`
			Name               string `json:"name"`
			UID                string `json:"uid"`
			Controller         bool   `json:"controller"`
			BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
		} `json:"ownerReferences"`
	} `json:"metadata"`
	Spec struct {
		Replicas int `json:"replicas"`
		Selector struct {
			MatchLabels map[string]string `json:"matchLabels"`
		} `json:"selector"`
		Template struct {
			Metadata struct {
				Name              string            `json:"name"`
				CreationTimestamp interface{}       `json:"creationTimestamp"`
				Labels            map[string]string `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				Containers []struct {
					Name  string `json:"name"`
					Image string `json:"image"`
					Ports []struct {
						Name          string `json:"name"`
						ContainerPort int    `json:"containerPort"`
						Protocol      string `json:"protocol"`
					} `json:"ports"`
					Resources struct {
					} `json:"resources"`
					TerminationMessagePath   string `json:"terminationMessagePath"`
					TerminationMessagePolicy string `json:"terminationMessagePolicy"`
					ImagePullPolicy          string `json:"imagePullPolicy"`
				} `json:"containers"`
				RestartPolicy                 string `json:"restartPolicy"`
				TerminationGracePeriodSeconds int    `json:"terminationGracePeriodSeconds"`
				DNSPolicy                     string `json:"dnsPolicy"`
				SecurityContext struct {
				} `json:"securityContext"`
				SchedulerName string `json:"schedulerName"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
	Status struct {
		Replicas             int `json:"replicas"`
		FullyLabeledReplicas int `json:"fullyLabeledReplicas"`
		ObservedGeneration   int `json:"observedGeneration"`
	} `json:"status"`
}
