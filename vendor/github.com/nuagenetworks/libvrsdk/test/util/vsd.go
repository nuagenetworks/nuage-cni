package util

import (
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

// VerifyVSDPortResolution will verify if the given port is present on VSD. If yes, returns the IP
func VerifyVSDPortResolution(root *vspk.Me, vsdEnterprise string, vsdDomain string, vsdZone string, vsdPort string) (string, *bambou.Error) {

	var enterprise *vspk.Enterprise
	var domain *vspk.Domain
	var zone *vspk.Zone
	var ipAddress string
	var err *bambou.Error

	//Fetching enterprise object
	enterprise, err = FetchEnterprise(root, vsdEnterprise)
	if err != nil {
		return ipAddress, err
	}

	//Fetching domain object
	domain, err = FetchDomain(enterprise, vsdDomain)
	if err != nil {
		return ipAddress, err
	}

	//Fetching zone object
	zone, err = FetchZone(domain, vsdZone)
	if err != nil {
		return ipAddress, err
	}

	//Fetching all port objects
	var ports vspk.VPortsList
	ports, err = FetchAllPorts(zone)
	if err != nil {
		return ipAddress, err
	}

	for i := 0; i < len(ports); i++ {

		vmInterfacesFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdPort + "\""}
		var vmInterfaces vspk.VMInterfacesList
		vmInterfaces, err = ports[i].VMInterfaces(vmInterfacesFetchingInfo)

		if err == nil && len(vmInterfaces) > 0 {
			vmInterface := vmInterfaces[0]
			ipAddress = vmInterface.IPAddress
		}
	}
	return ipAddress, err
}

// VerifyVSDPortDeletion will verify if the given port is removed from VSD or not
func VerifyVSDPortDeletion(root *vspk.Me, vsdEnterprise string, vsdDomain string, vsdZone string, vsdPort string) (bool, *bambou.Error) {

	var enterprise *vspk.Enterprise
	var domain *vspk.Domain
	var zone *vspk.Zone
	var portDeletionFailure bool
	var err *bambou.Error

	// Fetching enterprise object
	enterprise, err = FetchEnterprise(root, vsdEnterprise)
	if err != nil {
		return portDeletionFailure, err
	}

	// Fetching domain object
	domain, err = FetchDomain(enterprise, vsdDomain)
	if err != nil {
		return portDeletionFailure, err
	}

	// Fetching zone object
	zone, err = FetchZone(domain, vsdZone)
	if err != nil {
		return portDeletionFailure, err
	}

	// Fetching all port objects
	var ports vspk.VPortsList
	ports, err = FetchAllPorts(zone)
	if err != nil {
		return portDeletionFailure, err
	}

	for i := 0; i < len(ports); i++ {

		vmInterfacesFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdPort + "\""}
		var vmInterfaces vspk.VMInterfacesList
		vmInterfaces, err = ports[i].VMInterfaces(vmInterfacesFetchingInfo)

		if err != nil || len(vmInterfaces) > 0 {
			portDeletionFailure = true
		}
	}
	return portDeletionFailure, err
}

// FetchEnterprise fetches enterprise object
func FetchEnterprise(root *vspk.Me, vsdEnterprise string) (*vspk.Enterprise, *bambou.Error) {

	var enterprise *vspk.Enterprise
	enterpriseFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdEnterprise + "\""}
	enterprises, enterpriseErr := root.Enterprises(enterpriseFetchingInfo)
	if enterpriseErr == nil {
		enterprise = enterprises[0]
	}
	return enterprise, enterpriseErr
}

// FetchDomain fetches domain object
func FetchDomain(enterprise *vspk.Enterprise, vsdDomain string) (*vspk.Domain, *bambou.Error) {

	var domain *vspk.Domain
	domainFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdDomain + "\""}
	domains, domainErr := enterprise.Domains(domainFetchingInfo)
	if domainErr == nil {
		domain = domains[0]
	}
	return domain, domainErr
}

// FetchZone fetches zone object
func FetchZone(domain *vspk.Domain, vsdZone string) (*vspk.Zone, *bambou.Error) {

	var zone *vspk.Zone
	zoneFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdZone + "\""}
	zones, zonesErr := domain.Zones(zoneFetchingInfo)
	if zonesErr == nil {
		zone = zones[0]
	}
	return zone, zonesErr
}

// FetchAllPorts fetches all port objects
func FetchAllPorts(zone *vspk.Zone) (vspk.VPortsList, *bambou.Error) {

	portsFetchingInfo := &bambou.FetchingInfo{Filter: "type == \"VM\""}
	return zone.VPorts(portsFetchingInfo)
}
