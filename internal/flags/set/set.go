package set

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"endobit.io/stack/internal/flags"
)

type (
	flag[T comparable] struct {
		name       string
		acceptZero bool
		value      T
	}

	Appliance   struct{ flag[string] }
	Arch        struct{ flag[string] }
	Cluster     struct{ flag[string] }
	Environment struct{ flag[string] }
	Host        struct{ flag[string] }
	JSON        struct{ flag[bool] }
	Location    struct{ flag[string] }
	Make        struct{ flag[string] }
	Model       struct{ flag[string] }
	Rack        struct{ flag[string] }
	Rank        struct{ flag[uint32] }
	Rename      struct{ flag[string] }
	Slot        struct{ flag[uint32] }
	Template    struct{ flag[string] }
	TimeZone    struct{ flag[string] }
	HostType    struct{ flag[string] }
	Value       struct{ flag[string] }
	Zone        struct{ flag[string] }
)

func (f *flag[T]) AcceptZero(flags *pflag.FlagSet) {
	if flags.Changed(f.name) {
		f.acceptZero = true
	}
}

func (f flag[T]) IsSet() bool {
	var zero T

	return f.value != zero || f.acceptZero
}

func (f flag[T]) Val() T {
	return f.value
}

func (f flag[T]) Ptr() *T {
	var zero T

	if f.value == zero && !f.acceptZero {
		return nil
	}

	return &f.value
}

func (j *JSON) Add(fs *pflag.FlagSet, object string) {
	j.name = flags.JSON
	addBool(fs, &j.value, j.name, "output "+object+" as JSON")
}

func (a *Appliance) Add(fs *pflag.FlagSet, object string, req bool) {
	a.name = flags.Appliance
	addString(fs, &a.value, a.name, "appliance for the "+object, req)
}

func (a *Arch) Add(fs *pflag.FlagSet, object string) {
	a.name = flags.Arch
	addString(fs, &a.value, a.name, "architecture for the "+object, false)
}

func (c *Cluster) Add(fs *pflag.FlagSet, object string, req bool) {
	c.name = flags.Cluster
	addString(fs, &c.value, c.name, "cluster for the "+object, req)
}

func (e *Environment) Add(fs *pflag.FlagSet, object string, req bool) {
	e.name = flags.Environment
	addString(fs, &e.value, e.name, "environment for the "+object, req)
}

func (h *Host) Add(fs *pflag.FlagSet, object string, req bool) {
	h.name = flags.Host
	addString(fs, &h.value, h.name, "host for the "+object, req)
}

func (h *HostType) Add(fs *pflag.FlagSet, object string) {
	h.name = flags.HostType
	addString(fs, &h.value, h.name, "type for the "+object, false)
}

func (l *Location) Add(fs *pflag.FlagSet, object string) {
	l.name = flags.Location
	addString(fs, &l.value, l.name, "location for the "+object, false)
}

func (m *Make) Add(fs *pflag.FlagSet, object string, req bool) {
	m.name = flags.Make
	addString(fs, &m.value, m.name, "make for the "+object, req)
}

func (m *Model) Add(fs *pflag.FlagSet, object string, req bool) {
	m.name = flags.Model
	addString(fs, &m.value, m.name, "model for the "+object, req)
}

func (r *Rack) Add(fs *pflag.FlagSet, object string, req bool) {
	r.name = flags.Rack
	addString(fs, &r.value, r.name, "rack for the "+object, req)
}

func (r *Rank) Add(fs *pflag.FlagSet, object string) {
	r.name = flags.Rank
	addUint32(fs, &r.value, r.name, "rank for the "+object, false)
}

func (s *Slot) Add(fs *pflag.FlagSet, object string) {
	s.name = flags.Slot
	addUint32(fs, &s.value, s.name, "slot for the "+object, false)
}

func (r *Rename) Add(fs *pflag.FlagSet, object string) {
	r.name = flags.Rename
	addString(fs, &r.value, r.name, "rename the "+object, false)
}

func (t *Template) Add(fs *pflag.FlagSet, object string) {
	t.name = flags.Template
	addString(fs, &t.value, t.name, "template for the "+object, false)
}

func (t *TimeZone) Add(fs *pflag.FlagSet, object string) {
	t.name = flags.TimeZone
	addString(fs, &t.value, t.name, "time zone for the "+object, false)
}

func (v *Value) Add(fs *pflag.FlagSet, object string) {
	v.name = flags.Value
	addString(fs, &v.value, v.name, "value of the "+object, false)
}

func (z *Zone) Add(fs *pflag.FlagSet, object string, req bool) {
	z.name = flags.Zone

	def := os.Getenv("METAL_ZONE")
	fs.StringVar(&z.value, z.name, def, "zone for the "+object)

	if req && def == "" {
		must(cobra.MarkFlagRequired(fs, z.name))
	}

}

func addString(fs *pflag.FlagSet, store *string, name, usage string, req bool) {
	fs.StringVar(store, name, "", usage)
	if req {
		must(cobra.MarkFlagRequired(fs, name))
	}
}

func addUint32(fs *pflag.FlagSet, store *uint32, name, usage string, req bool) {
	fs.Uint32Var(store, name, 0, usage)
	if req {
		must(cobra.MarkFlagRequired(fs, name))
	}
}

func addBool(fs *pflag.FlagSet, store *bool, name, usage string) {
	fs.BoolVar(store, name, false, usage)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
