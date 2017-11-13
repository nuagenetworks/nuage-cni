package config

// NuageVSPK8SConfig struct will be used to read and
// parse values from Nuage vsp-k8s yaml file on k8s agent nodes
type NuageVSPK8SConfig struct {
	ClientCertFile            string `yaml:"clientCert"`
	ClientKeyFile             string `yaml:"clientKey"`
	CACertFile                string `yaml:"CACert"`
	EnterpriseName            string `yaml:"enterpriseName"`
	DomainName                string `yaml:"domainName"`
	VSDUser                   string `yaml:"vsdUser"`
	K8SAPIServer              string `yaml:"masterApiServer"`
	NuageK8SMonServer         string `yaml:"nuageMonRestServer"`
	DockerBridgeName          string `yaml:"dockerBridgeName"`
	ServiceCIDR               string `yaml:"serviceCIDR"`
	NuageK8SMonClientCertFile string `yaml:"nuageMonClientCert"`
	NuageK8SMonClientKeyFile  string `yaml:"nuageMonClientKey"`
	NuageK8SMonCAFile         string `yaml:"nuageMonCA"`
}

// Config struct will be used to read values from Nuage CNI
// parameter file necessary for audit daemon and CNI plugin
type Config struct {
	VRSEndpoint             string
	VRSBridge               string
	MonitorInterval         int
	CNIVersion              string
	LogLevel                string
	PortResolveTimer        int
	LogFileSize             int
	VRSConnectionCheckTimer int
	MTU                     int
	StaleEntryTimeout       int64
}
