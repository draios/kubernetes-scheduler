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

type KubeDeployments struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata struct {
	} `json:"metadata"`
	Items []KubeDeploymentItem `json:"items"`
}

type KubeDeploymentItem struct {
	Metadata struct {
		Name              string            `json:"name"`
		Namespace         string            `json:"namespace"`
		SelfLink          string            `json:"selfLink"`
		UID               string            `json:"uid"`
		ResourceVersion   string            `json:"resourceVersion"`
		Generation        int               `json:"generation"`
		CreationTimestamp time.Time         `json:"creationTimestamp"`
		Labels            map[string]string `json:"labels"`
		Annotations struct {
			DeploymentKubernetesIoRevision string `json:"deployment.kubernetes.io/revision"`
		} `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		Replicas int `json:"replicas"`
		Selector struct {
			MatchLabels map[string]string `json:"matchLabels"`
		} `json:"selector"`
		Template struct {
			Metadata struct {
				CreationTimestamp interface{}       `json:"creationTimestamp"`
				Labels            map[string]string `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				Containers []struct {
					Name  string `json:"name"`
					Image string `json:"image"`
					Resources struct {
						Limits struct {
							CPU    string `json:"cpu,omitempty"`
							Memory string `json:"memory,omitempty"`
						} `json:"limits,omitempty"`
						Requests struct {
							Memory string `json:"memory"`
						} `json:"requests"`
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
				SchedulerName string `json:"schedulerName,omitempty"`
			} `json:"spec"`
		} `json:"template"`
		Strategy struct {
			Type string `json:"type"`
			RollingUpdate struct {
				MaxUnavailable int `json:"maxUnavailable"`
				MaxSurge       int `json:"maxSurge"`
			} `json:"rollingUpdate"`
		} `json:"strategy"`
		RevisionHistoryLimit    int `json:"revisionHistoryLimit"`
		ProgressDeadlineSeconds int `json:"progressDeadlineSeconds"`
	} `json:"spec"`
	Status struct {
		ObservedGeneration  int `json:"observedGeneration"`
		Replicas            int `json:"replicas"`
		UpdatedReplicas     int `json:"updatedReplicas"`
		UnavailableReplicas int `json:"unavailableReplicas"`
		Conditions []struct {
			Type               string    `json:"type"`
			Status             string    `json:"status"`
			LastUpdateTime     time.Time `json:"lastUpdateTime"`
			LastTransitionTime time.Time `json:"lastTransitionTime"`
			Reason             string    `json:"reason"`
			Message            string    `json:"message"`
		} `json:"conditions"`
	} `json:"status"`
}
