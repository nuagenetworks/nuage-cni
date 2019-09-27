package k8s

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/nuagenetworks/nuage-cni/client"
	"github.com/nuagenetworks/nuage-cni/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var vspK8SConfig = &config.NuageVSPK8SConfig{}
var podNetwork string
var podZone string
var podPG string
var adminUser string
var podUID string

var vspK8sConfigFile string
var kubeconfFile string
var nuageMonClientCertFile string
var nuageMonClientKeyFile string
var nuageMonClientCACertFile string

var isHostAtomic bool

// NuageKubeMonResp will unmarshal JSON
// response from Nuage kubemon service
type NuageKubeMonResp struct {
	Subnet string   `json:"subnetName"`
	PG     []string `json:"policyGroups"`
}

// Pod will hold fields necessary to query
// Nuage kubemon service to obtain pod metadata
type Pod struct {
	Name   string `json:"podName"`
	Zone   string `json:"desiredZone,omitempty"`
	Subnet string `json:"desiredSubnet,omitempty"`
	Action string `json:"action,omitempty"`
}

func getK8SLabelsPodUIDFromAPIServer(podNs string, podname string) error {

	log.Infof("Obtaining labels from API server for pod %s under namespace %s", podname, podNs)

	// creates the in-cluster config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfFile)
	if err != nil {
		log.Errorf("failed to load kubeconfig %v", err)
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pod, err := clientset.CoreV1().Pods(podNs).Get(podname, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Error occured while querying pod %s under pod namespace %s: %v", podname, podNs, err)
		return err
	}

	podUID = string(pod.UID)
	if _, ok := pod.Labels["nuage.io/subnet"]; !ok {
		podNetwork = ""
	} else {
		podNetwork = pod.Labels["nuage.io/subnet"]
	}

	if _, ok := pod.Labels["nuage.io/zone"]; !ok {
		podZone = podNs
	} else {
		podZone = pod.Labels["nuage.io/zone"]
	}

	if _, ok := pod.Labels["nuage.io/user"]; !ok {
		adminUser = vspK8SConfig.VSDUser
	} else {
		adminUser = pod.Labels["nuage.io/user"]
	}

	if _, ok := pod.Labels["nuage.io/policy-group"]; !ok {
		podPG = ""
	} else {
		podPG = pod.Labels["nuage.io/policy-group"]
	}

	return err
}

func getVSPK8SConfig() error {

	// Reading Nuage VSP K8S yaml file
	data, err := ioutil.ReadFile(vspK8sConfigFile)
	if err != nil {
		return fmt.Errorf("Error in reading from Nuage VSP k8s yaml file: %s", err)
	}

	if err = yaml.Unmarshal(data, vspK8SConfig); err != nil {
		return fmt.Errorf("Error in unmarshalling data from Nuage VSP k8s yaml file: %s", err)
	}

	if vspK8SConfig.EnterpriseName == "" {
		vspK8SConfig.EnterpriseName = "K8S-Enterprise"
	}

	if vspK8SConfig.DomainName == "" {
		vspK8SConfig.DomainName = "K8S-Domain"
	}

	return err
}

func getPodMetadataFromNuageK8sMon(podname string, ns string) error {

	log.Infof("Obtaining Nuage Metadata for pod %s under namespace %s", podname, ns)
	var result = new(NuageKubeMonResp)
	url := vspK8SConfig.NuageK8SMonServer + "/namespaces/" + ns + "/pods"

	// Load client cert
	cert, err := tls.LoadX509KeyPair(nuageMonClientCertFile, nuageMonClientKeyFile)
	if err != nil {
		log.Errorf("Error loading client cert file to communicate with Nuage K8S monitor: %v", err)
		return err
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(nuageMonClientCACertFile)
	if err != nil {
		log.Errorf("Error loading CA cert file to communicate with Nuage K8S monitor: %v", err)
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	pod := &Pod{Name: podname}
	if podZone != ns {
		log.Infof("Desired zone %s and network %s set as labels for pod %s", podZone, podNetwork, podname)
		pod = &Pod{Name: podname, Zone: podZone, Subnet: podNetwork}
	}
	out, err := json.Marshal(pod)
	if err != nil {
		log.Errorf("Error occured while marshalling Pod data to communicate with Nuage K8S monitor: %v", err)
		return err
	}

	var jsonStr = []byte(string(out))
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Errorf("Error occured while sending POST call to Nuage K8S monitor to obtain pod metadata: %v", err)
		return err
	}

	log.Debugf("Response sent to Nuage kubemon is %v", bytes.NewBuffer(jsonStr))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error occured while reading response obtained from Nuage K8S monitor: %v", err)
		return err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Errorf("Error occured while unmarshalling Pod data obtained from Nuage K8S monitor: %v", err)
		return err
	}

	log.Debugf("Response status obtained from Nuage Kubemon for pod %s: %s", podname, resp.Status)
	log.Debugf("Result obtained as a result of passed labels for pod %s: %v", podname, result)

	if podPG == "" {
		log.Debugf("Pod policy group information obtained from Nuage K8S monitor : %s", result.PG)
		for _, pg := range result.PG {
			podPG = pg
		}
	}

	log.Debugf("Pod subnet information obtained from Nuage K8S monitor : %s", result.Subnet)
	podNetwork = result.Subnet

	return err
}

func initDataDir(orchestrator string) {

	isHostAtomic = VerifyHostType()
	var dir string
	if isHostAtomic == true {
		dir = "/var/usr/share/"
	} else {
		dir = "/usr/share/"
	}

	if orchestrator == "k8s" {
		vspK8sConfigFile = dir + "/vsp-k8s/vsp-k8s.yaml"
	} else {
		vspK8sConfigFile = dir + "/vsp-openshift/vsp-openshift.yaml"
		kubeconfFile = dir + "/vsp-openshift/nuage.kubeconfig"
		nuageMonClientCertFile = dir + "/vsp-openshift/nuageMonClient.crt"
		nuageMonClientKeyFile = dir + "/vsp-openshift/nuageMonClient.key"
		nuageMonClientCACertFile = dir + "/vsp-openshift/nuageMonCA.crt"
	}
}

// VerifyHostType will determine the base host
// as RHEL server or RHEL atomic
func VerifyHostType() bool {

	// check if the host is an atomic host
	_, err := os.Stat("/run/ostree-booted")
	if err != nil {
		log.Infof("This is a RHEL server host")
		return false
	}
	log.Infof("This is a RHEL atomic host")
	return true
}

// GetPodNuageMetadata will populate NuageMetadata struct
// needed for port resolution using CNI plugin
func GetPodNuageMetadata(nuageMetadata *client.NuageMetadata, name string, ns string, orchestrator string) error {

	initDataDir(orchestrator)
	log.Infof("Obtaining Nuage Metadata for pod %s under namespace %s", name, ns)

	var err error

	// Parsing Nuage VSP K8S yaml file on K8S agent nodes
	err = getVSPK8SConfig()
	if err != nil {
		log.Errorf("Error in parsing Nuage k8s yaml file")
		return fmt.Errorf("Error in parsing Nuage k8s yaml file: %s", err)
	}

	// Populating certificate and kubeconfig locations
	// only for k8s as orchestrator
	if orchestrator == "k8s" {
		kubeconfFile = vspK8SConfig.KubeConfig
		nuageMonClientCertFile = vspK8SConfig.NuageK8SMonClientCertFile
		nuageMonClientKeyFile = vspK8SConfig.NuageK8SMonClientKeyFile
		nuageMonClientCACertFile = vspK8SConfig.NuageK8SMonCAFile
	}

	// Obtaining pod labels if set from K8S API server
	err = getK8SLabelsPodUIDFromAPIServer(ns, name)
	if err != nil {
		log.Errorf("Error in obtaining pod labels from API server")
		return fmt.Errorf("Error in obtaining pod labels from API server: %s", err)
	}

	// Obtaining pod subnet/policy group metadata from Nuage K8S monitor service
	err = getPodMetadataFromNuageK8sMon(name, ns)
	if err != nil {
		log.Errorf("Error in obtaining pod subnet/policy group from Nuage K8S monitor")
		return fmt.Errorf("Error in obtaining pod subnet/policy group from Nuage K8S monitor: %s", err)
	}

	nuageMetadata.Enterprise = vspK8SConfig.EnterpriseName
	nuageMetadata.Domain = vspK8SConfig.DomainName
	nuageMetadata.Zone = podZone
	nuageMetadata.Network = podNetwork
	nuageMetadata.User = adminUser
	nuageMetadata.PolicyGroup = podPG
	nuageMetadata.PodUID = podUID

	return err
}

// SendPodDeletionNotification will notify the Nuage monitor on master nodes
// about pod deletion
func SendPodDeletionNotification(podname string, ns string, orchestrator string) error {

	initDataDir(orchestrator)
	var err error

	log.Infof("Sending delete notification for pod %s under namespace %s", podname, ns)

	// Parsing Nuage config file on agent nodes
	err = getVSPK8SConfig()
	if err != nil {
		log.Errorf("Error in parsing Nuage config file")
		return fmt.Errorf("Error in parsing Nuage config file: %s", err)
	}

	// Populating certificate and kubeconfig locations
	// only for k8s as orchestrator
	if orchestrator == "k8s" {
		kubeconfFile = vspK8SConfig.KubeConfig
		nuageMonClientCertFile = vspK8SConfig.NuageK8SMonClientCertFile
		nuageMonClientKeyFile = vspK8SConfig.NuageK8SMonClientKeyFile
		nuageMonClientCACertFile = vspK8SConfig.NuageK8SMonCAFile
	}

	url := vspK8SConfig.NuageK8SMonServer + "/namespaces/" + ns + "/pods"

	// Load client cert
	cert, err := tls.X509KeyPair([]byte(nuageMonClientCertFile), []byte(nuageMonClientKeyFile))
	if err != nil {
		log.Errorf("Error loading client cert file to communicate with Nuage monitor: %v", err)
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(nuageMonClientCACertFile))

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	pod := &Pod{Name: podname, Action: "delete"}
	out, err := json.Marshal(pod)
	if err != nil {
		log.Errorf("Error occured while marshalling Pod data to communicate with Nuage monitor: %v", err)
		return err
	}

	var jsonStr = []byte(string(out))
	_, err = client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Errorf("Error occured while sending pod deletion notification to Nuage monitor: %v", err)
		return err
	}

	return err
}
