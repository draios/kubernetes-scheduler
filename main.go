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

// A custom Kubernetes scheduler with custom metrics from Sysdig Monitor
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/draios/kubernetes-scheduler/cache"
	kube "github.com/draios/kubernetes-scheduler/kubernetes"
	"github.com/draios/kubernetes-scheduler/sysdig"
	"os/user"
	"time"
)

// Variables that will be used in our scheduler
var (
	schedulerName  string
	kubeAPI        kube.KubernetesCoreV1Api
	sysdigAPI      sysdig.SysdigApiClient
	metrics        []map[string]interface{}
	sysdigMetric   string
	bestCachedNode = cache.Cache{Timeout: 15 * time.Second}
	cachedNodes    = cache.Cache{Timeout: 15 * time.Second}
)

// Errors
var (
	noDataFound   = errors.New("no data found with those parameters")
	emptyNodeList = errors.New("node list must contain at least one element")
	noNodeFound   = errors.New("no node found")
)

// Flags
var (
	sysdigTokenFlag    = flag.String("t", "", "Sysdig Cloud Token")
	kubeConfigFileFlag = flag.String("k", "", "Kubernetes config file")
	sysdigMetricFlag   = flag.String("m", "", "Sysdig metric to monitorize")
	schedulerNameFlag  = flag.String("s", "", "Scheduler name")
)

func init() {
	usr, _ := user.Current()

	flag.Usage = usage
	flag.Parse()

	if sysdigTokenEnv, tokenSetByEnv := os.LookupEnv("SDC_TOKEN"); !tokenSetByEnv && *sysdigTokenFlag == "" {
		fmt.Println("Error: Sysdig Cloud token is not set.")
		usage()
	} else {
		if tokenSetByEnv {
			sysdigAPI.SetToken(sysdigTokenEnv)
		}
		if *sysdigTokenFlag != "" { // If the flag is set, overrides the environment
			sysdigAPI.SetToken(*sysdigTokenFlag)
		}
	}

	if _, kubeTokenSetByEnv := os.LookupEnv("KUBECONFIG"); !kubeTokenSetByEnv && *kubeConfigFileFlag == "" {
		os.Setenv("KUBECONFIG", usr.HomeDir+"/.kube/config")
	} else {
		if *kubeConfigFileFlag != "" {
			os.Setenv("KUBECONFIG", *kubeConfigFileFlag)
		}
	}
	kubeAPI.LoadKubeConfig()

	if sysdigMetricEnv, sysdigMetricEnvIsSet := os.LookupEnv("SDC_METRIC"); !sysdigMetricEnvIsSet && *sysdigMetricFlag == "" {
		fmt.Println("The Sysdig metric must be defined")
		usage()
	} else {
		if sysdigMetricEnvIsSet {
			sysdigMetric = sysdigMetricEnv
		}
		if *sysdigMetricFlag != "" {
			sysdigMetric = *sysdigMetricFlag
		}
	}

	if schedulerNameEnv, schedulernameEnvIsSet := os.LookupEnv("SDC_SCHEDULER"); !schedulernameEnvIsSet && *schedulerNameFlag == "" {
		fmt.Println("Scheduler name must be set")
		usage()
	} else {
		if schedulernameEnvIsSet {
			schedulerName = schedulerNameEnv
		}
		if *schedulerNameFlag != "" {
			schedulerName = *schedulerNameFlag
		}
	}
}

// Usage description
func usage() {
	fmt.Printf("Usage: %s [-s SCHEDULER_NAME] [-m SYSDIG_METRIC] [-t SYSDIG_TOKEN] [-k KUBERNETES_CONFIG_FILE]", os.Args[0])
	fmt.Print(`
If the env KUBECONFIG is not set, the -k option must be provided.
If the env SDC_TOKEN is not set, the -t option must be provided.
If the env SDC_METRIC is not set, the -m option must be provided.
If the env SDC_SCHEDULER is not set, the -s option must be provided.
`)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	kubePod := KubePod{}

	metrics = append(metrics, map[string]interface{}{
		"id": sysdigMetric,
		"aggregations": map[string]string{
			"time": "timeAvg", "group": "avg",
		},
	})

	ch, _ := kubeAPI.Watch("GET", "api/v1/pods", nil, nil)
	for data := range ch {
		err := json.Unmarshal(data, &kubePod)
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		if kubePod.Object.Status.Phase == "Pending" && kubePod.Object.Spec.SchedulerName == schedulerName && kubePod.Type == "ADDED" {
			log.Println("Scheduling", kubePod.Object.Metadata.Name)

			bestNodeFound, err := getBestRequestTime(nodesAvailable())
			if err != nil {
				log.Println("Error:", err)
			} else {
				log.Println("Best node found: ", bestNodeFound.name, bestNodeFound.time)
				response, err := scheduler(kubePod.Object.Metadata.Name, bestNodeFound.name, kubePod.Object.Metadata.Namespace)
				if err != nil {
					log.Println("Error:", err)
				}
				kubeResponse := kube.KubeResponse{}
				err = json.NewDecoder(response.Body).Decode(&kubeResponse)
				if err != nil {
					log.Println("error while decoding kube response: ", err)
				}
				if kubeResponse.Code != 200 && kubeResponse.Code != 201 {
					log.Println("kube response error: ", kubeResponse.Message)
				}

				response.Body.Close()
			}
		}
	}
}
