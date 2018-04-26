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

type KubeResponse struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata struct {
	} `json:"metadata"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Details struct {
		Name string `json:"name"`
		Kind string `json:"kind"`
	} `json:"details"`
	Code int `json:"code"`
}
