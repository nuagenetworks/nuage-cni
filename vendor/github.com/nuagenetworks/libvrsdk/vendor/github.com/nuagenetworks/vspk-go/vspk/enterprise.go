/*
  Copyright (c) 2015, Alcatel-Lucent Inc
  All rights reserved.

  Redistribution and use in source and binary forms, with or without
  modification, are permitted provided that the following conditions are met:
      * Redistributions of source code must retain the above copyright
        notice, this list of conditions and the following disclaimer.
      * Redistributions in binary form must reproduce the above copyright
        notice, this list of conditions and the following disclaimer in the
        documentation and/or other materials provided with the distribution.
      * Neither the name of the copyright holder nor the names of its contributors
        may be used to endorse or promote products derived from this software without
        specific prior written permission.

  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
  ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
  WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
  DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY
  DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
  (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
  LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
  ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
  (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
  SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package vspk

import "github.com/nuagenetworks/go-bambou/bambou"

// EnterpriseIdentity represents the Identity of the object
var EnterpriseIdentity = bambou.Identity{
	Name:     "enterprise",
	Category: "enterprises",
}

// EnterprisesList represents a list of Enterprises
type EnterprisesList []*Enterprise

// EnterprisesAncestor is the interface of an ancestor of a Enterprise must implement.
type EnterprisesAncestor interface {
	Enterprises(*bambou.FetchingInfo) (EnterprisesList, *bambou.Error)
	CreateEnterprises(*Enterprise) *bambou.Error
}

// Enterprise represents the model of a enterprise
type Enterprise struct {
	ID                                    string        `json:"ID,omitempty"`
	ParentID                              string        `json:"parentID,omitempty"`
	ParentType                            string        `json:"parentType,omitempty"`
	Owner                                 string        `json:"owner,omitempty"`
	LDAPAuthorizationEnabled              bool          `json:"LDAPAuthorizationEnabled"`
	LDAPEnabled                           bool          `json:"LDAPEnabled"`
	BGPEnabled                            bool          `json:"BGPEnabled"`
	DHCPLeaseInterval                     int           `json:"DHCPLeaseInterval,omitempty"`
	Name                                  string        `json:"name,omitempty"`
	LastUpdatedBy                         string        `json:"lastUpdatedBy,omitempty"`
	ReceiveMultiCastListID                string        `json:"receiveMultiCastListID,omitempty"`
	SendMultiCastListID                   string        `json:"sendMultiCastListID,omitempty"`
	Description                           string        `json:"description,omitempty"`
	AllowAdvancedQOSConfiguration         bool          `json:"allowAdvancedQOSConfiguration"`
	AllowGatewayManagement                bool          `json:"allowGatewayManagement"`
	AllowTrustedForwardingClass           bool          `json:"allowTrustedForwardingClass"`
	AllowedForwardingClasses              []interface{} `json:"allowedForwardingClasses,omitempty"`
	FloatingIPsQuota                      int           `json:"floatingIPsQuota,omitempty"`
	FloatingIPsUsed                       int           `json:"floatingIPsUsed,omitempty"`
	EncryptionManagementMode              string        `json:"encryptionManagementMode,omitempty"`
	EnterpriseProfileID                   string        `json:"enterpriseProfileID,omitempty"`
	EntityScope                           string        `json:"entityScope,omitempty"`
	LocalAS                               int           `json:"localAS,omitempty"`
	AssociatedEnterpriseSecurityID        string        `json:"associatedEnterpriseSecurityID,omitempty"`
	AssociatedGroupKeyEncryptionProfileID string        `json:"associatedGroupKeyEncryptionProfileID,omitempty"`
	AssociatedKeyServerMonitorID          string        `json:"associatedKeyServerMonitorID,omitempty"`
	CustomerID                            int           `json:"customerID,omitempty"`
	AvatarData                            string        `json:"avatarData,omitempty"`
	AvatarType                            string        `json:"avatarType,omitempty"`
	ExternalID                            string        `json:"externalID,omitempty"`
}

// NewEnterprise returns a new *Enterprise
func NewEnterprise() *Enterprise {

	return &Enterprise{}
}

// Identity returns the Identity of the object.
func (o *Enterprise) Identity() bambou.Identity {

	return EnterpriseIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Enterprise) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Enterprise) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Enterprise from the server
func (o *Enterprise) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Enterprise into the server
func (o *Enterprise) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Enterprise from the server
func (o *Enterprise) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// L2Domains retrieves the list of child L2Domains of the Enterprise
func (o *Enterprise) L2Domains(info *bambou.FetchingInfo) (L2DomainsList, *bambou.Error) {

	var list L2DomainsList
	err := bambou.CurrentSession().FetchChildren(o, L2DomainIdentity, &list, info)
	return list, err
}

// CreateL2Domain creates a new child L2Domain under the Enterprise
func (o *Enterprise) CreateL2Domain(child *L2Domain) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// L2DomainTemplates retrieves the list of child L2DomainTemplates of the Enterprise
func (o *Enterprise) L2DomainTemplates(info *bambou.FetchingInfo) (L2DomainTemplatesList, *bambou.Error) {

	var list L2DomainTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, L2DomainTemplateIdentity, &list, info)
	return list, err
}

// CreateL2DomainTemplate creates a new child L2DomainTemplate under the Enterprise
func (o *Enterprise) CreateL2DomainTemplate(child *L2DomainTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// RateLimiters retrieves the list of child RateLimiters of the Enterprise
func (o *Enterprise) RateLimiters(info *bambou.FetchingInfo) (RateLimitersList, *bambou.Error) {

	var list RateLimitersList
	err := bambou.CurrentSession().FetchChildren(o, RateLimiterIdentity, &list, info)
	return list, err
}

// CreateRateLimiter creates a new child RateLimiter under the Enterprise
func (o *Enterprise) CreateRateLimiter(child *RateLimiter) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Gateways retrieves the list of child Gateways of the Enterprise
func (o *Enterprise) Gateways(info *bambou.FetchingInfo) (GatewaysList, *bambou.Error) {

	var list GatewaysList
	err := bambou.CurrentSession().FetchChildren(o, GatewayIdentity, &list, info)
	return list, err
}

// CreateGateway creates a new child Gateway under the Enterprise
func (o *Enterprise) CreateGateway(child *Gateway) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GatewayTemplates retrieves the list of child GatewayTemplates of the Enterprise
func (o *Enterprise) GatewayTemplates(info *bambou.FetchingInfo) (GatewayTemplatesList, *bambou.Error) {

	var list GatewayTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, GatewayTemplateIdentity, &list, info)
	return list, err
}

// CreateGatewayTemplate creates a new child GatewayTemplate under the Enterprise
func (o *Enterprise) CreateGatewayTemplate(child *GatewayTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PATNATPools retrieves the list of child PATNATPools of the Enterprise
func (o *Enterprise) PATNATPools(info *bambou.FetchingInfo) (PATNATPoolsList, *bambou.Error) {

	var list PATNATPoolsList
	err := bambou.CurrentSession().FetchChildren(o, PATNATPoolIdentity, &list, info)
	return list, err
}

// CreatePATNATPool creates a new child PATNATPool under the Enterprise
func (o *Enterprise) CreatePATNATPool(child *PATNATPool) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// LDAPConfigurations retrieves the list of child LDAPConfigurations of the Enterprise
func (o *Enterprise) LDAPConfigurations(info *bambou.FetchingInfo) (LDAPConfigurationsList, *bambou.Error) {

	var list LDAPConfigurationsList
	err := bambou.CurrentSession().FetchChildren(o, LDAPConfigurationIdentity, &list, info)
	return list, err
}

// CreateLDAPConfiguration creates a new child LDAPConfiguration under the Enterprise
func (o *Enterprise) CreateLDAPConfiguration(child *LDAPConfiguration) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// RedundancyGroups retrieves the list of child RedundancyGroups of the Enterprise
func (o *Enterprise) RedundancyGroups(info *bambou.FetchingInfo) (RedundancyGroupsList, *bambou.Error) {

	var list RedundancyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, RedundancyGroupIdentity, &list, info)
	return list, err
}

// CreateRedundancyGroup creates a new child RedundancyGroup under the Enterprise
func (o *Enterprise) CreateRedundancyGroup(child *RedundancyGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the Enterprise
func (o *Enterprise) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Enterprise
func (o *Enterprise) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// MetadataTags retrieves the list of child MetadataTags of the Enterprise
func (o *Enterprise) MetadataTags(info *bambou.FetchingInfo) (MetadataTagsList, *bambou.Error) {

	var list MetadataTagsList
	err := bambou.CurrentSession().FetchChildren(o, MetadataTagIdentity, &list, info)
	return list, err
}

// CreateMetadataTag creates a new child MetadataTag under the Enterprise
func (o *Enterprise) CreateMetadataTag(child *MetadataTag) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NetworkMacroGroups retrieves the list of child NetworkMacroGroups of the Enterprise
func (o *Enterprise) NetworkMacroGroups(info *bambou.FetchingInfo) (NetworkMacroGroupsList, *bambou.Error) {

	var list NetworkMacroGroupsList
	err := bambou.CurrentSession().FetchChildren(o, NetworkMacroGroupIdentity, &list, info)
	return list, err
}

// CreateNetworkMacroGroup creates a new child NetworkMacroGroup under the Enterprise
func (o *Enterprise) CreateNetworkMacroGroup(child *NetworkMacroGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// KeyServerMonitors retrieves the list of child KeyServerMonitors of the Enterprise
func (o *Enterprise) KeyServerMonitors(info *bambou.FetchingInfo) (KeyServerMonitorsList, *bambou.Error) {

	var list KeyServerMonitorsList
	err := bambou.CurrentSession().FetchChildren(o, KeyServerMonitorIdentity, &list, info)
	return list, err
}

// CreateKeyServerMonitor creates a new child KeyServerMonitor under the Enterprise
func (o *Enterprise) CreateKeyServerMonitor(child *KeyServerMonitor) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ZFBRequests retrieves the list of child ZFBRequests of the Enterprise
func (o *Enterprise) ZFBRequests(info *bambou.FetchingInfo) (ZFBRequestsList, *bambou.Error) {

	var list ZFBRequestsList
	err := bambou.CurrentSession().FetchChildren(o, ZFBRequestIdentity, &list, info)
	return list, err
}

// CreateZFBRequest creates a new child ZFBRequest under the Enterprise
func (o *Enterprise) CreateZFBRequest(child *ZFBRequest) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// BGPProfiles retrieves the list of child BGPProfiles of the Enterprise
func (o *Enterprise) BGPProfiles(info *bambou.FetchingInfo) (BGPProfilesList, *bambou.Error) {

	var list BGPProfilesList
	err := bambou.CurrentSession().FetchChildren(o, BGPProfileIdentity, &list, info)
	return list, err
}

// CreateBGPProfile creates a new child BGPProfile under the Enterprise
func (o *Enterprise) CreateBGPProfile(child *BGPProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressQOSPolicies retrieves the list of child EgressQOSPolicies of the Enterprise
func (o *Enterprise) EgressQOSPolicies(info *bambou.FetchingInfo) (EgressQOSPoliciesList, *bambou.Error) {

	var list EgressQOSPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, EgressQOSPolicyIdentity, &list, info)
	return list, err
}

// CreateEgressQOSPolicy creates a new child EgressQOSPolicy under the Enterprise
func (o *Enterprise) CreateEgressQOSPolicy(child *EgressQOSPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// SharedNetworkResources retrieves the list of child SharedNetworkResources of the Enterprise
func (o *Enterprise) SharedNetworkResources(info *bambou.FetchingInfo) (SharedNetworkResourcesList, *bambou.Error) {

	var list SharedNetworkResourcesList
	err := bambou.CurrentSession().FetchChildren(o, SharedNetworkResourceIdentity, &list, info)
	return list, err
}

// CreateSharedNetworkResource creates a new child SharedNetworkResource under the Enterprise
func (o *Enterprise) CreateSharedNetworkResource(child *SharedNetworkResource) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKECertificates retrieves the list of child IKECertificates of the Enterprise
func (o *Enterprise) IKECertificates(info *bambou.FetchingInfo) (IKECertificatesList, *bambou.Error) {

	var list IKECertificatesList
	err := bambou.CurrentSession().FetchChildren(o, IKECertificateIdentity, &list, info)
	return list, err
}

// CreateIKECertificate creates a new child IKECertificate under the Enterprise
func (o *Enterprise) CreateIKECertificate(child *IKECertificate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEEncryptionprofiles retrieves the list of child IKEEncryptionprofiles of the Enterprise
func (o *Enterprise) IKEEncryptionprofiles(info *bambou.FetchingInfo) (IKEEncryptionprofilesList, *bambou.Error) {

	var list IKEEncryptionprofilesList
	err := bambou.CurrentSession().FetchChildren(o, IKEEncryptionprofileIdentity, &list, info)
	return list, err
}

// CreateIKEEncryptionprofile creates a new child IKEEncryptionprofile under the Enterprise
func (o *Enterprise) CreateIKEEncryptionprofile(child *IKEEncryptionprofile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEGateways retrieves the list of child IKEGateways of the Enterprise
func (o *Enterprise) IKEGateways(info *bambou.FetchingInfo) (IKEGatewaysList, *bambou.Error) {

	var list IKEGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, IKEGatewayIdentity, &list, info)
	return list, err
}

// CreateIKEGateway creates a new child IKEGateway under the Enterprise
func (o *Enterprise) CreateIKEGateway(child *IKEGateway) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEGatewayProfiles retrieves the list of child IKEGatewayProfiles of the Enterprise
func (o *Enterprise) IKEGatewayProfiles(info *bambou.FetchingInfo) (IKEGatewayProfilesList, *bambou.Error) {

	var list IKEGatewayProfilesList
	err := bambou.CurrentSession().FetchChildren(o, IKEGatewayProfileIdentity, &list, info)
	return list, err
}

// CreateIKEGatewayProfile creates a new child IKEGatewayProfile under the Enterprise
func (o *Enterprise) CreateIKEGatewayProfile(child *IKEGatewayProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEPSKs retrieves the list of child IKEPSKs of the Enterprise
func (o *Enterprise) IKEPSKs(info *bambou.FetchingInfo) (IKEPSKsList, *bambou.Error) {

	var list IKEPSKsList
	err := bambou.CurrentSession().FetchChildren(o, IKEPSKIdentity, &list, info)
	return list, err
}

// CreateIKEPSK creates a new child IKEPSK under the Enterprise
func (o *Enterprise) CreateIKEPSK(child *IKEPSK) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the Enterprise
func (o *Enterprise) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// CreateAlarm creates a new child Alarm under the Enterprise
func (o *Enterprise) CreateAlarm(child *Alarm) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// AllAlarms retrieves the list of child AllAlarms of the Enterprise
func (o *Enterprise) AllAlarms(info *bambou.FetchingInfo) (AllAlarmsList, *bambou.Error) {

	var list AllAlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AllAlarmIdentity, &list, info)
	return list, err
}

// CreateAllAlarm creates a new child AllAlarm under the Enterprise
func (o *Enterprise) CreateAllAlarm(child *AllAlarm) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Enterprise
func (o *Enterprise) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Enterprise
func (o *Enterprise) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the Enterprise
func (o *Enterprise) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// CreateVM creates a new child VM under the Enterprise
func (o *Enterprise) CreateVM(child *VM) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// InfrastructurePortProfiles retrieves the list of child InfrastructurePortProfiles of the Enterprise
func (o *Enterprise) InfrastructurePortProfiles(info *bambou.FetchingInfo) (InfrastructurePortProfilesList, *bambou.Error) {

	var list InfrastructurePortProfilesList
	err := bambou.CurrentSession().FetchChildren(o, InfrastructurePortProfileIdentity, &list, info)
	return list, err
}

// CreateInfrastructurePortProfile creates a new child InfrastructurePortProfile under the Enterprise
func (o *Enterprise) CreateInfrastructurePortProfile(child *InfrastructurePortProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterpriseNetworks retrieves the list of child EnterpriseNetworks of the Enterprise
func (o *Enterprise) EnterpriseNetworks(info *bambou.FetchingInfo) (EnterpriseNetworksList, *bambou.Error) {

	var list EnterpriseNetworksList
	err := bambou.CurrentSession().FetchChildren(o, EnterpriseNetworkIdentity, &list, info)
	return list, err
}

// CreateEnterpriseNetwork creates a new child EnterpriseNetwork under the Enterprise
func (o *Enterprise) CreateEnterpriseNetwork(child *EnterpriseNetwork) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterpriseSecurities retrieves the list of child EnterpriseSecurities of the Enterprise
func (o *Enterprise) EnterpriseSecurities(info *bambou.FetchingInfo) (EnterpriseSecuritiesList, *bambou.Error) {

	var list EnterpriseSecuritiesList
	err := bambou.CurrentSession().FetchChildren(o, EnterpriseSecurityIdentity, &list, info)
	return list, err
}

// CreateEnterpriseSecurity creates a new child EnterpriseSecurity under the Enterprise
func (o *Enterprise) CreateEnterpriseSecurity(child *EnterpriseSecurity) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Jobs retrieves the list of child Jobs of the Enterprise
func (o *Enterprise) Jobs(info *bambou.FetchingInfo) (JobsList, *bambou.Error) {

	var list JobsList
	err := bambou.CurrentSession().FetchChildren(o, JobIdentity, &list, info)
	return list, err
}

// CreateJob creates a new child Job under the Enterprise
func (o *Enterprise) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Domains retrieves the list of child Domains of the Enterprise
func (o *Enterprise) Domains(info *bambou.FetchingInfo) (DomainsList, *bambou.Error) {

	var list DomainsList
	err := bambou.CurrentSession().FetchChildren(o, DomainIdentity, &list, info)
	return list, err
}

// CreateDomain creates a new child Domain under the Enterprise
func (o *Enterprise) CreateDomain(child *Domain) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DomainTemplates retrieves the list of child DomainTemplates of the Enterprise
func (o *Enterprise) DomainTemplates(info *bambou.FetchingInfo) (DomainTemplatesList, *bambou.Error) {

	var list DomainTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, DomainTemplateIdentity, &list, info)
	return list, err
}

// CreateDomainTemplate creates a new child DomainTemplate under the Enterprise
func (o *Enterprise) CreateDomainTemplate(child *DomainTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Containers retrieves the list of child Containers of the Enterprise
func (o *Enterprise) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// CreateContainer creates a new child Container under the Enterprise
func (o *Enterprise) CreateContainer(child *Container) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// RoutingPolicies retrieves the list of child RoutingPolicies of the Enterprise
func (o *Enterprise) RoutingPolicies(info *bambou.FetchingInfo) (RoutingPoliciesList, *bambou.Error) {

	var list RoutingPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, RoutingPolicyIdentity, &list, info)
	return list, err
}

// CreateRoutingPolicy creates a new child RoutingPolicy under the Enterprise
func (o *Enterprise) CreateRoutingPolicy(child *RoutingPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Apps retrieves the list of child Apps of the Enterprise
func (o *Enterprise) Apps(info *bambou.FetchingInfo) (AppsList, *bambou.Error) {

	var list AppsList
	err := bambou.CurrentSession().FetchChildren(o, AppIdentity, &list, info)
	return list, err
}

// CreateApp creates a new child App under the Enterprise
func (o *Enterprise) CreateApp(child *App) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ApplicationServices retrieves the list of child ApplicationServices of the Enterprise
func (o *Enterprise) ApplicationServices(info *bambou.FetchingInfo) (ApplicationServicesList, *bambou.Error) {

	var list ApplicationServicesList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationServiceIdentity, &list, info)
	return list, err
}

// CreateApplicationService creates a new child ApplicationService under the Enterprise
func (o *Enterprise) CreateApplicationService(child *ApplicationService) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Groups retrieves the list of child Groups of the Enterprise
func (o *Enterprise) Groups(info *bambou.FetchingInfo) (GroupsList, *bambou.Error) {

	var list GroupsList
	err := bambou.CurrentSession().FetchChildren(o, GroupIdentity, &list, info)
	return list, err
}

// CreateGroup creates a new child Group under the Enterprise
func (o *Enterprise) CreateGroup(child *Group) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GroupKeyEncryptionProfiles retrieves the list of child GroupKeyEncryptionProfiles of the Enterprise
func (o *Enterprise) GroupKeyEncryptionProfiles(info *bambou.FetchingInfo) (GroupKeyEncryptionProfilesList, *bambou.Error) {

	var list GroupKeyEncryptionProfilesList
	err := bambou.CurrentSession().FetchChildren(o, GroupKeyEncryptionProfileIdentity, &list, info)
	return list, err
}

// CreateGroupKeyEncryptionProfile creates a new child GroupKeyEncryptionProfile under the Enterprise
func (o *Enterprise) CreateGroupKeyEncryptionProfile(child *GroupKeyEncryptionProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DSCPForwardingClassTables retrieves the list of child DSCPForwardingClassTables of the Enterprise
func (o *Enterprise) DSCPForwardingClassTables(info *bambou.FetchingInfo) (DSCPForwardingClassTablesList, *bambou.Error) {

	var list DSCPForwardingClassTablesList
	err := bambou.CurrentSession().FetchChildren(o, DSCPForwardingClassTableIdentity, &list, info)
	return list, err
}

// CreateDSCPForwardingClassTable creates a new child DSCPForwardingClassTable under the Enterprise
func (o *Enterprise) CreateDSCPForwardingClassTable(child *DSCPForwardingClassTable) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Users retrieves the list of child Users of the Enterprise
func (o *Enterprise) Users(info *bambou.FetchingInfo) (UsersList, *bambou.Error) {

	var list UsersList
	err := bambou.CurrentSession().FetchChildren(o, UserIdentity, &list, info)
	return list, err
}

// CreateUser creates a new child User under the Enterprise
func (o *Enterprise) CreateUser(child *User) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGateways retrieves the list of child NSGateways of the Enterprise
func (o *Enterprise) NSGateways(info *bambou.FetchingInfo) (NSGatewaysList, *bambou.Error) {

	var list NSGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewayIdentity, &list, info)
	return list, err
}

// CreateNSGateway creates a new child NSGateway under the Enterprise
func (o *Enterprise) CreateNSGateway(child *NSGateway) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGatewayTemplates retrieves the list of child NSGatewayTemplates of the Enterprise
func (o *Enterprise) NSGatewayTemplates(info *bambou.FetchingInfo) (NSGatewayTemplatesList, *bambou.Error) {

	var list NSGatewayTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewayTemplateIdentity, &list, info)
	return list, err
}

// CreateNSGatewayTemplate creates a new child NSGatewayTemplate under the Enterprise
func (o *Enterprise) CreateNSGatewayTemplate(child *NSGatewayTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSRedundantGatewayGroups retrieves the list of child NSRedundantGatewayGroups of the Enterprise
func (o *Enterprise) NSRedundantGatewayGroups(info *bambou.FetchingInfo) (NSRedundantGatewayGroupsList, *bambou.Error) {

	var list NSRedundantGatewayGroupsList
	err := bambou.CurrentSession().FetchChildren(o, NSRedundantGatewayGroupIdentity, &list, info)
	return list, err
}

// CreateNSRedundantGatewayGroup creates a new child NSRedundantGatewayGroup under the Enterprise
func (o *Enterprise) CreateNSRedundantGatewayGroup(child *NSRedundantGatewayGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PublicNetworkMacros retrieves the list of child PublicNetworkMacros of the Enterprise
func (o *Enterprise) PublicNetworkMacros(info *bambou.FetchingInfo) (PublicNetworkMacrosList, *bambou.Error) {

	var list PublicNetworkMacrosList
	err := bambou.CurrentSession().FetchChildren(o, PublicNetworkMacroIdentity, &list, info)
	return list, err
}

// CreatePublicNetworkMacro creates a new child PublicNetworkMacro under the Enterprise
func (o *Enterprise) CreatePublicNetworkMacro(child *PublicNetworkMacro) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// MultiCastLists retrieves the list of child MultiCastLists of the Enterprise
func (o *Enterprise) MultiCastLists(info *bambou.FetchingInfo) (MultiCastListsList, *bambou.Error) {

	var list MultiCastListsList
	err := bambou.CurrentSession().FetchChildren(o, MultiCastListIdentity, &list, info)
	return list, err
}

// CreateMultiCastList creates a new child MultiCastList under the Enterprise
func (o *Enterprise) CreateMultiCastList(child *MultiCastList) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Avatars retrieves the list of child Avatars of the Enterprise
func (o *Enterprise) Avatars(info *bambou.FetchingInfo) (AvatarsList, *bambou.Error) {

	var list AvatarsList
	err := bambou.CurrentSession().FetchChildren(o, AvatarIdentity, &list, info)
	return list, err
}

// CreateAvatar creates a new child Avatar under the Enterprise
func (o *Enterprise) CreateAvatar(child *Avatar) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the Enterprise
func (o *Enterprise) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}

// CreateEventLog creates a new child EventLog under the Enterprise
func (o *Enterprise) CreateEventLog(child *EventLog) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ExternalAppServices retrieves the list of child ExternalAppServices of the Enterprise
func (o *Enterprise) ExternalAppServices(info *bambou.FetchingInfo) (ExternalAppServicesList, *bambou.Error) {

	var list ExternalAppServicesList
	err := bambou.CurrentSession().FetchChildren(o, ExternalAppServiceIdentity, &list, info)
	return list, err
}

// CreateExternalAppService creates a new child ExternalAppService under the Enterprise
func (o *Enterprise) CreateExternalAppService(child *ExternalAppService) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ExternalServices retrieves the list of child ExternalServices of the Enterprise
func (o *Enterprise) ExternalServices(info *bambou.FetchingInfo) (ExternalServicesList, *bambou.Error) {

	var list ExternalServicesList
	err := bambou.CurrentSession().FetchChildren(o, ExternalServiceIdentity, &list, info)
	return list, err
}

// CreateExternalService creates a new child ExternalService under the Enterprise
func (o *Enterprise) CreateExternalService(child *ExternalService) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
