package kubernetes

import "time"

type KubePodEvent struct {
	Type string `json:"type"`
	// Object of the event
	Object struct {
		Kind       string `json:"kind"`
		APIVersion string `json:"apiVersion"`
		Metadata struct {
			Name              string            `json:"name"`
			GenerateName      string            `json:"generateName"`
			Namespace         string            `json:"namespace"`
			SelfLink          string            `json:"selfLink"`
			UID               string            `json:"uid"`
			ResourceVersion   string            `json:"resourceVersion"`
			CreationTimestamp time.Time         `json:"creationTimestamp"`
			Labels            map[string]string `json:"labels"`
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
			Volumes []struct {
				Name string `json:"name"`
				Secret struct {
					SecretName  string `json:"secretName"`
					DefaultMode int    `json:"defaultMode"`
				} `json:"secret"`
			} `json:"volumes"`
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
				VolumeMounts []struct {
					Name      string `json:"name"`
					ReadOnly  bool   `json:"readOnly"`
					MountPath string `json:"mountPath"`
				} `json:"volumeMounts"`
				TerminationMessagePath   string `json:"terminationMessagePath"`
				TerminationMessagePolicy string `json:"terminationMessagePolicy"`
				ImagePullPolicy          string `json:"imagePullPolicy"`
			} `json:"containers"`
			RestartPolicy                 string `json:"restartPolicy"`
			TerminationGracePeriodSeconds int    `json:"terminationGracePeriodSeconds"`
			DNSPolicy                     string `json:"dnsPolicy"`
			ServiceAccountName            string `json:"serviceAccountName"`
			ServiceAccount                string `json:"serviceAccount"`
			NodeName                      string `json:"nodeName"`
			SecurityContext struct {
			} `json:"securityContext"`
			SchedulerName string `json:"schedulerName"`
			Tolerations []struct {
				Key               string `json:"key"`
				Operator          string `json:"operator"`
				Effect            string `json:"effect"`
				TolerationSeconds int    `json:"tolerationSeconds"`
			} `json:"tolerations"`
		} `json:"spec"`
		Status struct {
			Phase string `json:"phase"`
			Conditions []struct {
				Type               string      `json:"type"`
				Status             string      `json:"status"`
				LastProbeTime      interface{} `json:"lastProbeTime"`
				LastTransitionTime time.Time   `json:"lastTransitionTime"`
			} `json:"conditions"`
			HostIP    string    `json:"hostIP"`
			PodIP     string    `json:"podIP"`
			StartTime time.Time `json:"startTime"`
			ContainerStatuses []struct {
				Name string `json:"name"`
				State struct {
					Running struct {
						StartedAt time.Time `json:"startedAt"`
					} `json:"running"`
				} `json:"state"`
				LastState struct {
				} `json:"lastState"`
				Ready        bool   `json:"ready"`
				RestartCount int    `json:"restartCount"`
				Image        string `json:"image"`
				ImageID      string `json:"imageID"`
				ContainerID  string `json:"containerID"`
			} `json:"containerStatuses"`
			QosClass string `json:"qosClass"`
		} `json:"status"`
	} `json:"object"`
}
