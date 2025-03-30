package commands

import "errors"

//go:generate gotip tool "github.com/dmarkham/enumer" -type Verb -transform lower -text

type Verb int

const (
	Add Verb = iota
	Dump
	List
	Load
	Remove
	Report
	Set
	Unset
)

const (
	attribute   = "attr"
	rack        = "rack"
	appliance   = "appliance"
	cluster     = "cluster"
	host        = "host"
	environment = "environment"
	model       = "model"
	zone        = "zone"
)

var (
	errInvalidHostType    = errors.New("invalid host type")
	errMissingClusterZone = errors.New("cluster zone not specified")
	errMissingMakeOrModel = errors.New("if either make or model is specified, both must be set")
)
