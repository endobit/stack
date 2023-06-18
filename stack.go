// Package stack provides a data model for a stack of hardware.
package stack

import "net/netip"

// Document is the top level data structure for a stack of hardware.
type Document struct {
	Models       []Model       `json:"models,omitempty" yaml:"models,omitempty"`
	Appliances   []Appliance   `json:"appliances,omitempty" yaml:"appliances,omitempty"`
	Environments []Environment `json:"environments,omitempty" yaml:"environments,omitempty"`
	Zones        []Zone        `json:"zones,omitempty" yaml:"zones,omitempty"`
}

// Attribute is a key/value pair with an optional protection flag.
type Attribute struct {
	Key       string `json:"key" yaml:"key"`
	Value     string `json:"value" yaml:"value"`
	Protected bool   `json:"protected,omitempty" yaml:"protected,omitempty"`
}

// Zone is a logical grouping of hardware at a physical location.
type Zone struct {
	Name       string      `json:"name" yaml:"name"`
	TimeZone   string      `json:"time_zone" yaml:"time_zone"`
	Attributes []Attribute `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	Networks   []Network   `json:"networks,omitempty" yaml:"networks,omitempty"`
	Clusters   []Cluster   `json:"clusters,omitempty" yaml:"clusters,omitempty"`
	Hosts      []Host      `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Switches   []Switch    `json:"switches,omitempty" yaml:"switches,omitempty"`
}

// Cluster is a tightly coupled group of hosts and switches.
type Cluster struct {
	Name       string      `json:"name" yaml:"name"`
	Zone       string      `json:"zone" yaml:"zone"`
	Attributes []Attribute `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	Networks   []Network   `json:"networks,omitempty" yaml:"networks,omitempty"`
	Hosts      []Host      `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Switches   []Switch    `json:"switches,omitempty" yaml:"switches,omitempty"`
}

// Appliance is a logical role for a host.
type Appliance struct {
	Name       string      `json:"name" yaml:"name"`
	Attributes []Attribute `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

// Model is a physical type of hardware.
type Model struct {
	Name       string      `json:"name" yaml:"name"`
	Attributes []Attribute `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

// Environment is a logical grouping of hosts or switches.
type Environment struct {
	Name       string      `json:"name" yaml:"name"`
	Attributes []Attribute `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

// Network is a network subnet.
type Network struct {
	Name    string      `json:"name" yaml:"name"`
	Address *netip.Addr `json:"address,omitempty" yaml:"address,omitempty"`
	Gateway *netip.Addr `json:"gateway,omitempty" yaml:"gateway,omitempty"`
	MTU     int         `json:"mtu,omitempty" yaml:"mtu,omitempty"`
	PXE     bool        `json:"pxe,omitempty" yaml:"pxe,omitempty"`
}

// Host is a physical or virtual server.
type Host struct {
	Name        string             `json:"name" yaml:"name"`
	Attributes  []Attribute        `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	Appliance   string             `json:"appliance,omitempty" yaml:"appliance,omitempty"`
	Environment string             `json:"environment,omitempty" yaml:"environment,omitempty"`
	Model       string             `json:"model,omitempty" yaml:"model,omitempty"`
	Rack        string             `json:"rack,omitempty" yaml:"rack,omitempty"`
	Rank        int                `json:"rank,omitempty" yaml:"rank,omitempty"`
	NICs        []NetworkInterface `json:"nics,omitempty" yaml:"nics,omitempty"`
	Type        HostType           `json:"type,omitempty" yaml:"type,omitempty"`
}

// Switch is a network switch.
type Switch struct {
	Name        string             `json:"name" yaml:"name"`
	Attributes  []Attribute        `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	Appliance   string             `json:"appliance,omitempty" yaml:"appliance,omitempty"`
	Environment string             `json:"environment,omitempty" yaml:"environment,omitempty"`
	Model       string             `json:"model,omitempty" yaml:"model,omitempty"`
	Rack        string             `json:"rack,omitempty" yaml:"rack,omitempty"`
	Rank        int                `json:"rank,omitempty" yaml:"rank,omitempty"`
	NICs        []NetworkInterface `json:"nics,omitempty" yaml:"nics,omitempty"`
}

// NetworkInterface is a network interface on a host or switch.
type NetworkInterface struct {
	IP         *netip.Addr `json:"ip,omitempty" yaml:"ip,omitempty"`
	MAC        string      `json:"mac" yaml:"mac"`
	DHCP       bool        `json:"dhcp,omitempty" yaml:"dhcp,omitempty"`
	Management bool        `json:"management,omitempty" yaml:"management,omitempty"`
}
