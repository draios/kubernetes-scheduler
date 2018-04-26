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

// Small Kubernetes api wrapper
package kubernetes

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"

	"gopkg.in/yaml.v2"
	"github.com/draios/kubernetes-scheduler/cache"
	"time"
	"net"
	"context"
)

type KubernetesCoreV1Api struct {
	config       KubeConf
	nodeList     cache.Cache
	clientCert   tls.Certificate
	serverCaCert *x509.CertPool
}

func getKubeConfigFileDefaultLocation() string {
	kubeConf, isSet := os.LookupEnv("KUBECONFIG")
	if isSet && kubeConf != "" {
		return kubeConf
	}

	usr, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	return usr.HomeDir + "/.kube/config"

}

func (api KubernetesCoreV1Api) CreateNamespacedBinding(namespace string, body io.Reader) (response *http.Response, err error) {
	return api.Request("POST", fmt.Sprintf("api/v1/namespaces/%s/bindings", namespace), nil, body)
}

func (api KubernetesCoreV1Api) Watch(httpMethod, apiMethod string, values url.Values, body io.Reader) (responseChannel chan []byte, err error) {
	if values == nil {
		values = url.Values{}
	}
	values.Add("watch", "true")
	responseChannel = make(chan []byte)
	go func() {
		response, err := api.Request(httpMethod, apiMethod, values, body)
		if err != nil {
			close(responseChannel)
			return
		}
		defer response.Body.Close()

		reader := bufio.NewReader(response.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				log.Println(err)
				continue
			}
			responseChannel <- line
			line = nil
		}
	}()
	return
}

func (api KubernetesCoreV1Api) Request(httpMethod, apiMethod string, values url.Values, body io.Reader) (response *http.Response, err error) {
	apiUrl := api.currentApiUrlEndpoint()

	certificate, caCertPool := api.currentTLSInfo()
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      caCertPool,
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig, DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, 1*time.Second)
	}}
	client := http.Client{Transport: transport}
	request, err := http.NewRequest(httpMethod, apiUrl+"/"+apiMethod, body)
	if err != nil {
		log.Println(err)
		return
	}
	if values != nil {
		request.URL.RawQuery = values.Encode()
	}

	// Get the info in json
	request.Header.Add("Content-Type", "application/json")

	// Make the request
	response, err = client.Do(request)
	if err != nil {
		log.Println(err)
	}
	return
}

func (api KubernetesCoreV1Api) ListNodes() (nodes []KubeNode, err error) {

	if nodes, ok := api.nodeList.Data(); ok {
		return nodes.([]KubeNode), nil
	}

	response, err := api.Request("GET", "api/v1/nodes", nil, nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	var nodeInfo struct {
		Items []KubeNode `json:"items"`
	}
	err = json.NewDecoder(response.Body).Decode(&nodeInfo)
	if err != nil {
		return
	}
	nodes = nodeInfo.Items

	api.nodeList.SetData(nodes)

	return
}

func (api KubernetesCoreV1Api) currentApiUrlEndpoint() string {
	for _, context := range api.config.Contexts {
		if context.Name == api.config.CurrentContext {
			for _, cluster := range api.config.Clusters {
				if cluster.Name == context.Data.Cluster {
					return cluster.Data.Server
				}
			}
		}
	}
	log.Panic("Current API Url endpoint couldn't be determined, checkout if the configuration is correct")
	return ""
}

// Parses the cert data and generates
func (api *KubernetesCoreV1Api) loadTLSInfo() {
	var currentContextUser string
	var currentContextCluster string
	var certData []byte
	var keyData []byte
	var caCertData []byte

	// Load current context information
	for _, context := range api.config.Contexts {
		if context.Name == api.config.CurrentContext {
			currentContextUser = context.Data.User
			currentContextCluster = context.Data.Cluster
		}
	}

	// Get cert and key data from the user
	for _, user := range api.config.Users {
		if user.Name == currentContextUser {
			certData, keyData = user.Data.ClientCertificateData, user.Data.ClientKeyData
		}
	}

	// Get CA Cert data from current cluster
	for _, cluster := range api.config.Clusters {
		if cluster.Name == currentContextCluster {
			caCertData = cluster.Data.CertificateAuthorityData
		}
	}

	certificate, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		panic(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertData)

	api.clientCert = certificate
	api.serverCaCert = caCertPool
}

func (api KubernetesCoreV1Api) currentTLSInfo() (clientCert tls.Certificate, serverCaCert *x509.CertPool) {
	return api.clientCert, api.serverCaCert
}

// Reads the configuration file and loads the config struct
func (api *KubernetesCoreV1Api) LoadKubeConfig() (err error) {
	yamlFile, err := ioutil.ReadFile(getKubeConfigFileDefaultLocation())
	if err != nil {
		panic("Could not load the Kubernetes configuration")
	}

	kubeConfig := KubeConf{}
	err = yaml.Unmarshal(yamlFile, &kubeConfig)
	if err != nil {
		return err
	}

	// Decode the certificate
	for k, cluster := range kubeConfig.Clusters {
		certBytes, err := base64.StdEncoding.DecodeString(cluster.Data.CertificateAuthorityDataStr)
		if err != nil {
			panic(err)
		}
		cluster.Data.CertificateAuthorityData = certBytes
		kubeConfig.Clusters[k] = cluster
	}

	// Decode the certificate and the key
	for k, user := range kubeConfig.Users {
		cert, err := base64.StdEncoding.DecodeString(user.Data.ClientCertificateDataStr)
		if err != nil {
			panic(err)
		}
		user.Data.ClientCertificateData = cert
		key, err := base64.StdEncoding.DecodeString(user.Data.ClientKeyDataStr)
		if err != nil {
			panic(err)
		}
		user.Data.ClientKeyData = key
		kubeConfig.Users[k] = user
	}

	api.config = kubeConfig
	api.nodeList = cache.Cache{Timeout: 1 * time.Minute}
	api.loadTLSInfo()
	return err
}
