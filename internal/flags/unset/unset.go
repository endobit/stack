package unset

import (
	"github.com/spf13/pflag"

	"endobit.io/stack/internal/flags"
)

type (
	flag struct {
		value bool
	}

	Appliance   struct{ flag }
	Arch        struct{ flag }
	Cluster     struct{ flag }
	Environment struct{ flag }
	HostType    struct{ flag }
	Location    struct{ flag }
	Make        struct{ flag }
	Model       struct{ flag }
	Rack        struct{ flag }
	Rank        struct{ flag }
	Slot        struct{ flag }
	Template    struct{ flag }
	TimeZone    struct{ flag }
	Value       struct{ flag }
)

func (f flag) Val() bool {
	return f.value
}

func (f flag) Ptr() *bool {
	if !f.value {
		return nil
	}

	return &f.value
}

func (a *Appliance) Add(fs *pflag.FlagSet, object string) {
	add(fs, &a.value, flags.Appliance, "appliance for the "+object)
}

func (a *Arch) Add(fs *pflag.FlagSet, object string) {
	add(fs, &a.value, flags.Arch, "architecture for the "+object)
}

func (c *Cluster) Add(fs *pflag.FlagSet, object string) {
	add(fs, &c.value, flags.Cluster, "cluster for the "+object)
}

func (e *Environment) Add(fs *pflag.FlagSet, object string) {
	add(fs, &e.value, flags.Environment, "environment for the "+object)
}

func (h *HostType) Add(fs *pflag.FlagSet, object string) {
	add(fs, &h.value, flags.HostType, "type for the "+object)
}

func (l *Location) Add(fs *pflag.FlagSet, object string) {
	add(fs, &l.value, flags.Location, "location for the "+object)
}

func (m *Make) Add(fs *pflag.FlagSet, object string) {
	add(fs, &m.value, flags.Make, "make for the "+object)
}

func (m *Model) Add(fs *pflag.FlagSet, object string) {
	add(fs, &m.value, flags.Model, "model for the "+object)
}

func (r *Rack) Add(fs *pflag.FlagSet, object string) {
	add(fs, &r.value, flags.Rack, "rack for the "+object)
}

func (r *Rank) Add(fs *pflag.FlagSet, object string) {
	add(fs, &r.value, flags.Rank, "rank for the "+object)
}

func (s *Slot) Add(fs *pflag.FlagSet, object string) {
	add(fs, &s.value, flags.Slot, "slot for the "+object)
}

func (t *TimeZone) Add(fs *pflag.FlagSet, object string) {
	add(fs, &t.value, flags.TimeZone, "time zone for the "+object)
}

func add(fs *pflag.FlagSet, store *bool, name, usage string) {
	fs.BoolVar(store, name, false, "unset the "+usage)
}
