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

type KubeConf struct {
	ApiVersion     string    `yaml:"apiVersion"`
	Clusters       []Cluster `yaml:"clusters"`
	Contexts       []Context `yaml:"contexts"`
	CurrentContext string    `yaml:"current-context"`
	Kind           string    `yaml:"kind"`
	Users          []User    `yaml:"users"`
}

type Cluster struct {
	Data ClusterData `yaml:"cluster"`
	Name string      `yaml:"name"`
}

type ClusterData struct {
	CertificateAuthorityDataStr string `yaml:"certificate-authority-data"`
	CertificateAuthorityData    []byte
	Server                      string `yaml:"server"`
}

type Context struct {
	Data ContextData `yaml:"context"`
	Name string      `yaml:"name"`
}

type ContextData struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type User struct {
	Name         string       `yaml:"name"`
	Data         UserData     `yaml:"user"`
	AuthProvider AuthProvider `yaml:"auth-provider"`
}

type UserData struct {
	ClientCertificateDataStr string `yaml:"client-certificate-data"`
	ClientCertificateData    []byte
	ClientKeyDataStr         string `yaml:"client-key-data"`
	ClientKeyData            []byte
}

type AuthProvider struct {
	Config AuthProviderConfig `yaml:"config"`
	Name   string             `yaml:"name"`
}

type AuthProviderConfig struct {
	CmdArgs   string `yaml:"cmd-args"`
	CmdPath   string `yaml:"cmd-path"`
	ExpiryKey string `yaml:"expiry-key"`
	TokenKey  string `yaml:"token-key"`
}
