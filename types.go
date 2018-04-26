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

package main

type KubePod struct {
	Type   string `json:"type"`
	Object struct {
		Status struct {
			Phase string `json:"phase"`
		} `json:"status"`
		Spec struct {
			SchedulerName string `json:"schedulerName"`
		} `json:"spec"`
		Metadata struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"metadata"`
	} `json:"object"`
}

type Node struct {
	name string
	time float64
	err  error
}
