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

// InfrastructurePortProfileIdentity represents the Identity of the object
var InfrastructurePortProfileIdentity = bambou.Identity{
	Name:     "infrastructureportprofile",
	Category: "infrastructureportprofiles",
}

// InfrastructurePortProfilesList represents a list of InfrastructurePortProfiles
type InfrastructurePortProfilesList []*InfrastructurePortProfile

// InfrastructurePortProfilesAncestor is the interface of an ancestor of a InfrastructurePortProfile must implement.
type InfrastructurePortProfilesAncestor interface {
	InfrastructurePortProfiles(*bambou.FetchingInfo) (InfrastructurePortProfilesList, *bambou.Error)
	CreateInfrastructurePortProfiles(*InfrastructurePortProfile) *bambou.Error
}

// InfrastructurePortProfile represents the model of a infrastructureportprofile
type InfrastructurePortProfile struct {
	ID            string `json:"ID,omitempty"`
	ParentID      string `json:"parentID,omitempty"`
	ParentType    string `json:"parentType,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Name          string `json:"name,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	Description   string `json:"description,omitempty"`
	EnterpriseID  string `json:"enterpriseID,omitempty"`
	EntityScope   string `json:"entityScope,omitempty"`
	Speed         string `json:"speed,omitempty"`
	UplinkTag     string `json:"uplinkTag,omitempty"`
	Mtu           int    `json:"mtu,omitempty"`
	Duplex        string `json:"duplex,omitempty"`
	ExternalID    string `json:"externalID,omitempty"`
}

// NewInfrastructurePortProfile returns a new *InfrastructurePortProfile
func NewInfrastructurePortProfile() *InfrastructurePortProfile {

	return &InfrastructurePortProfile{
		Duplex: "FULL",
		Speed:  "AUTONEGOTIATE",
		Mtu:    1500,
	}
}

// Identity returns the Identity of the object.
func (o *InfrastructurePortProfile) Identity() bambou.Identity {

	return InfrastructurePortProfileIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *InfrastructurePortProfile) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *InfrastructurePortProfile) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the InfrastructurePortProfile from the server
func (o *InfrastructurePortProfile) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the InfrastructurePortProfile into the server
func (o *InfrastructurePortProfile) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the InfrastructurePortProfile from the server
func (o *InfrastructurePortProfile) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the InfrastructurePortProfile
func (o *InfrastructurePortProfile) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the InfrastructurePortProfile
func (o *InfrastructurePortProfile) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the InfrastructurePortProfile
func (o *InfrastructurePortProfile) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the InfrastructurePortProfile
func (o *InfrastructurePortProfile) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
