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

// AppIdentity represents the Identity of the object
var AppIdentity = bambou.Identity{
	Name:     "application",
	Category: "applications",
}

// AppsList represents a list of Apps
type AppsList []*App

// AppsAncestor is the interface of an ancestor of a App must implement.
type AppsAncestor interface {
	Apps(*bambou.FetchingInfo) (AppsList, *bambou.Error)
	CreateApps(*App) *bambou.Error
}

// App represents the model of a application
type App struct {
	ID                          string `json:"ID,omitempty"`
	ParentID                    string `json:"parentID,omitempty"`
	ParentType                  string `json:"parentType,omitempty"`
	Owner                       string `json:"owner,omitempty"`
	Name                        string `json:"name,omitempty"`
	LastUpdatedBy               string `json:"lastUpdatedBy,omitempty"`
	Description                 string `json:"description,omitempty"`
	EntityScope                 string `json:"entityScope,omitempty"`
	AssocEgressACLTemplateId    string `json:"assocEgressACLTemplateId,omitempty"`
	AssocIngressACLTemplateId   string `json:"assocIngressACLTemplateId,omitempty"`
	AssociatedDomainID          string `json:"associatedDomainID,omitempty"`
	AssociatedDomainType        string `json:"associatedDomainType,omitempty"`
	AssociatedNetworkObjectID   string `json:"associatedNetworkObjectID,omitempty"`
	AssociatedNetworkObjectType string `json:"associatedNetworkObjectType,omitempty"`
	ExternalID                  string `json:"externalID,omitempty"`
}

// NewApp returns a new *App
func NewApp() *App {

	return &App{}
}

// Identity returns the Identity of the object.
func (o *App) Identity() bambou.Identity {

	return AppIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *App) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *App) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the App from the server
func (o *App) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the App into the server
func (o *App) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the App from the server
func (o *App) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the App
func (o *App) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the App
func (o *App) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Tiers retrieves the list of child Tiers of the App
func (o *App) Tiers(info *bambou.FetchingInfo) (TiersList, *bambou.Error) {

	var list TiersList
	err := bambou.CurrentSession().FetchChildren(o, TierIdentity, &list, info)
	return list, err
}

// CreateTier creates a new child Tier under the App
func (o *App) CreateTier(child *Tier) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the App
func (o *App) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the App
func (o *App) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Flows retrieves the list of child Flows of the App
func (o *App) Flows(info *bambou.FetchingInfo) (FlowsList, *bambou.Error) {

	var list FlowsList
	err := bambou.CurrentSession().FetchChildren(o, FlowIdentity, &list, info)
	return list, err
}

// CreateFlow creates a new child Flow under the App
func (o *App) CreateFlow(child *Flow) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Jobs retrieves the list of child Jobs of the App
func (o *App) Jobs(info *bambou.FetchingInfo) (JobsList, *bambou.Error) {

	var list JobsList
	err := bambou.CurrentSession().FetchChildren(o, JobIdentity, &list, info)
	return list, err
}

// CreateJob creates a new child Job under the App
func (o *App) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the App
func (o *App) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}

// CreateEventLog creates a new child EventLog under the App
func (o *App) CreateEventLog(child *EventLog) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
