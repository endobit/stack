// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"database/sql"
)

type Appliance struct {
	ID        int64
	Appliance string
}

type Cluster struct {
	ID      int64
	ZoneID  int64
	Cluster string
}

type ClusterAttribute struct {
	ID          int64
	ClusterID   int64
	Key         string
	Value       string
	IsProtected int64
}

type ClusterNetwork struct {
	NetworkID int64
	ClusterID int64
}

type Device struct {
	ID            int64
	Device        string
	ClusterID     sql.NullInt64
	ApplianceID   sql.NullInt64
	ModelID       sql.NullInt64
	EnvironmentID sql.NullInt64
}

type DeviceAttribute struct {
	ID          int64
	DeviceID    int64
	Key         string
	Value       string
	IsProtected int64
}

type Environment struct {
	ID          int64
	Environment string
}

type EnvironmentAttribute struct {
	ID            int64
	EnvironmentID int64
	Key           string
	Value         string
	IsProtected   int64
}

type GlobalAttribute struct {
	ID          int64
	Key         string
	Value       string
	IsProtected int64
}

type Host struct {
	ID       int64
	DeviceID sql.NullInt64
	Rack     sql.NullString
	Rank     sql.NullInt64
	Type     sql.NullString
}

type Model struct {
	ID    int64
	Model string
	Arch  string
}

type ModelAttribute struct {
	ID          int64
	ModelID     int64
	Key         string
	Value       string
	IsProtected int64
}

type Network struct {
	ID      int64
	Network string
	Address string
	Gateway sql.NullString
	Mtu     int64
	IsPxe   int64
}

type NetworkInterfac struct {
	ID           int64
	DeviceID     sql.NullInt64
	NetworkID    sql.NullInt64
	Ip           sql.NullString
	Mac          sql.NullString
	IsDhcp       int64
	IsManagement int64
}

type Switch struct {
	ID       int64
	DeviceID sql.NullInt64
	Rack     sql.NullString
	Rank     sql.NullInt64
}

type Zone struct {
	ID       int64
	Zone     string
	TimeZone string
}

type ZoneAttribute struct {
	ID          int64
	ZoneID      int64
	Key         string
	Value       string
	IsProtected int64
}

type ZoneNetwork struct {
	NetworkID int64
	ZoneID    int64
}
