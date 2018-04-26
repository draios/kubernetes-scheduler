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

// Small api wrapper to retrieve information from Sysdig Monitor
package sysdig

import (
	"net/http"
	"fmt"
	"encoding/json"
	"bytes"
	"io"
	"time"
)

const apiUrl = "https://api.sysdigcloud.com/"

type SysdigApiClient struct {
	token string
}

func (api *SysdigApiClient) SetToken(token string) {
	api.token = token
}

// Export metric data (both time-series and table-based)
//
// - metrics:
// 		A list of dictionaries, specifying the metrics and grouping keys that the query will return.
// 		A metric is any of the entries that can be found in the *Metrics* section of the Explore page in Sysdig Monitor.
// 		Metric entries require an *aggregations* section specifying how to aggregate the metric across time
// 		and containers/hosts.
// 		A grouping key is any of the entries that can be found in the *Show* or *Segment By* sections of
// 		the Explore page in Sysdig Monitor.
// 		These entries are used to apply single or hierarchical segmentation to the returned data and don't
// 		require the aggregations section.
//
// - start:
// 		The UTC time (in seconds) of the beginning of the data window.
// 		A negative value can be optionally used to indicate a relative time in the past from now.
// 		For example, -3600 means "one hour ago".
//
// - end:
// 		The UTC time (in seconds) of the end of the data window, or 0 to indicate "now".
// 		A negative value can also be optionally used to indicate a relative time in the past from now.
// 		For example, -3600 means "one hour ago".
//
// - sampling:
// 		The duration of the samples that will be returned.
// 		0 means that the whole data will be returned as a single sample.
//
// - filter:
// 		A boolean expression combining Sysdig Monitor segmentation criteria that defines what the query will be
// 		applied to.
// 		For example: *kubernetes.namespace.name='production' and container.image='nginx'*.
//
// - datasourceType:
// 		Specify the metric source for the request, can be "container" or "host".
// 		Most metrics, for example "cpu.used.percent" or "memory.bytes.used", are reported by both hosts and containers.
// 		By default, host metrics are used, but if the request contains a container-specific grouping key in the metric
// 		list/filter (e.g. "container.name"), then the container source is used.
// 		In cases where grouping keys are missing or apply to both hosts and containers (e.g. "tag.Name"),
// 		datasourceType can be explicitly set to avoid any ambiguity and allow the user to select precisely what kind of
// 		data should be used for the request.
func (api SysdigApiClient) GetData(metrics []map[string]interface{}, start, end, sampling int, filter, dataSourceType string) (response *http.Response, err error) {
	if dataSourceType == "" {
		dataSourceType = "host"
	}

	reqBody := map[string]interface{}{
		"metrics":        metrics,
		"dataSourceType": dataSourceType,
	}

	if start < 0 {
		reqBody["last"] = -start
	} else if start == 0 {
		err = fmt.Errorf("%s", "startTs cannot be 0")
		return
	} else {
		reqBody["start"] = start
		reqBody["end"] = end
	}

	if filter != "" {
		reqBody["filter"] = filter
	}

	if sampling != 0 {
		reqBody["sampling"] = sampling
	}

	reqBytes, err := json.Marshal(reqBody)
	body := bytes.NewReader(reqBytes)

	return api.Request("POST", "api/data", body)
}

// Makes a request to the Sysdig API endpoint.
//
// - httpMethod:
// 		The HTTP request method ("GET", "POST", "PUT", ...).
//
// - apiMethod:
// 		The API endpoint ("api/data", "api/token", ...).
//
// - body:
// 		Information that will be sent to the endpoint.
func (api SysdigApiClient) Request(httpMethod, apiMethod string, body io.Reader) (response *http.Response, err error) {

	// Create the request
	client := http.Client{Timeout: 5 * time.Second}
	request, err := http.NewRequest(httpMethod, apiUrl+apiMethod, body)
	if err != nil {
		return
	}
	// Header needed to connect with Sysdig Cloud
	request.Header.Add("Authorization", "Bearer "+api.token)
	// Get the info in json
	request.Header.Add("Content-Type", "application/json")

	// Make the request
	response, err = client.Do(request)
	return
}
