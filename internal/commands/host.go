package commands

import (
	"strconv"

	"github.com/spf13/cobra"

	"endobit.io/metal"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/stack/internal/flags/set"
	"endobit.io/stack/internal/flags/unset"
	"endobit.io/table"
)

type Host struct {
	*Root
	make        set.Make
	model       set.Model
	environment set.Environment
	appliance   set.Appliance
	location    set.Location
	rack        set.Rack
	rank        set.Rank
	slot        set.Slot
	hostType    set.HostType
}

func NewHost(r *Root) *Host {
	return &Host{Root: r}
}

type HostAttr struct {
	*Host
	host set.Host
}

func NewHostAttr(h *Host) *HostAttr {
	return &HostAttr{Host: h}
}

func (h *Host) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   host + " name",
		Short: "Add a " + host + " to a zone or cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := h.create(args[0]); err != nil {
				return err
			}

			return h.update(args[0])
		},
	}

	h.zone.Add(cmd.Flags(), host, true)
	h.cluster.Add(cmd.Flags(), host, false)

	cmd.AddCommand(NewHostAttr(h).Add())

	return cmd
}

func (h *Host) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   host + " name",
		Short: "Set a " + host + "'s properties",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(h.model.Val())+len(h.make.Val()) == 1 {
				return errMissingMakeOrModel
			}

			if h.hostType.IsSet() {
				if _, ok := pb.HostType_value[h.hostType.Val()]; !ok {
					return errInvalidHostType
				}
			}

			h.rank.AcceptZero(cmd.Flags())
			h.slot.AcceptZero(cmd.Flags())

			return h.update(args[0])
		},
	}

	h.zone.Add(cmd.Flags(), host, true)
	h.cluster.Add(cmd.Flags(), host, false)

	h.rename.Add(cmd.Flags(), host)
	h.make.Add(cmd.Flags(), host, false)
	h.model.Add(cmd.Flags(), host, false)
	h.environment.Add(cmd.Flags(), host, false)
	h.appliance.Add(cmd.Flags(), host, false)
	h.location.Add(cmd.Flags(), host)
	h.rack.Add(cmd.Flags(), host, false)
	h.rank.Add(cmd.Flags(), host)
	h.slot.Add(cmd.Flags(), host)
	h.hostType.Add(cmd.Flags(), host)

	cmd.AddCommand(NewHostAttr(h).Set())

	return cmd
}

func (h *Host) Unset() *cobra.Command {
	var (
		make        unset.Make
		model       unset.Model
		environment unset.Environment
		appliance   unset.Appliance
		location    unset.Location
		rack        unset.Rack
		rank        unset.Rank
		slot        unset.Slot
		hostType    unset.HostType
	)

	cmd := &cobra.Command{
		Use:   host + " name",
		Short: "Unset a " + host + "'s properties",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var m *bool

			if make.Val() || model.Val() {
				m = Ptr(true) // make/model are set together, so unset together
			}

			req := pb.UpdateHostRequest_builder{
				Zone:    h.zone.Ptr(),
				Cluster: h.cluster.Ptr(),
				Name:    &args[0],
				Unset: pb.UpdateHostRequest_Unset_builder{
					Make:        m,
					Model:       m,
					Environment: environment.Ptr(),
					Appliance:   appliance.Ptr(),
					Location:    location.Ptr(),
					Rack:        rack.Ptr(),
					Rank:        rank.Ptr(),
					Slot:        slot.Ptr(),
					Type:        hostType.Ptr(),
				}.Build(),
			}.Build()

			_, err := h.Metal.UpdateHost(h.Metal.Context(), req)

			return err

		},
	}

	h.zone.Add(cmd.Flags(), host, true)
	h.cluster.Add(cmd.Flags(), host, false)

	make.Add(cmd.Flags(), host)
	model.Add(cmd.Flags(), host)
	environment.Add(cmd.Flags(), host)
	appliance.Add(cmd.Flags(), host)
	location.Add(cmd.Flags(), host)
	rack.Add(cmd.Flags(), host)
	rank.Add(cmd.Flags(), host)
	slot.Add(cmd.Flags(), host)
	hostType.Add(cmd.Flags(), host)

	cmd.AddCommand(NewHostAttr(h).Set())

	return cmd
}

func (h *Host) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   host + " [glob]",
		Short: "List one or more " + host + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}

			return h.list(glob)
		},
	}

	h.zone.Add(cmd.Flags(), host, false)
	h.cluster.Add(cmd.Flags(), host, false)

	cmd.AddCommand(NewHostAttr(h).List())

	return cmd
}

func (h *Host) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   host + " glob",
		Short: "Remove one or more " + host + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return h.remove(args[0])
		},
	}

	h.zone.Add(cmd.Flags(), host, true)
	h.cluster.Add(cmd.Flags(), host, false)

	cmd.AddCommand(NewHostAttr(h).Remove())

	return cmd
}

func (h *Host) create(host string) error {
	req := pb.CreateHostRequest_builder{
		Zone:    h.zone.Ptr(),
		Cluster: h.cluster.Ptr(),
		Name:    &host,
	}.Build()

	_, err := h.Metal.CreateHost(h.Metal.Context(), req)

	return err
}

func (h *Host) list(glob string) error {
	type row struct {
		Zone        string
		Cluster     string `table:",omitempty"`
		Host        string
		Location    string `table:",omitempty"`
		Make, Model string
		Environment string `table:",omitempty"`
		Appliance   string
		Rack        string
		Rank        string `table:",omitempty"`
		Slot        string `table:",omitempty"`
		Type        string `table:",omitempty"`
	}

	t := table.New()
	defer t.Flush()

	r := h.Metal.NewHostReader(h.zone.Val(), h.cluster.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		var rank, slot, hostType string

		if resp.HasRank() {
			rank = strconv.Itoa(int(resp.GetRank()))
		}
		if resp.HasSlot() {
			slot = strconv.Itoa(int(resp.GetSlot()))
		}
		if resp.HasType() {
			hostType = metal.ShortHostType[resp.GetType().String()]
		}

		_ = t.Write(row{
			Zone:        resp.GetZone(),
			Cluster:     resp.GetCluster(),
			Host:        resp.GetName(),
			Make:        resp.GetMake(),
			Model:       resp.GetModel(),
			Environment: resp.GetEnvironment(),
			Appliance:   resp.GetAppliance(),
			Location:    resp.GetLocation(),
			Rack:        resp.GetRack(),
			Rank:        rank,
			Slot:        slot,
			Type:        hostType,
		})
	}

	return nil
}

func (h *Host) update(host string) error {
	var ht *pb.HostType

	if h.hostType.IsSet() {
		ht = Ptr(pb.HostType(pb.HostType_value[h.hostType.Val()]))
	}

	req := pb.UpdateHostRequest_builder{
		Zone:    h.zone.Ptr(),
		Cluster: h.cluster.Ptr(),
		Name:    &host,
		Set: pb.UpdateHostRequest_Set_builder{
			Name:        h.rename.Ptr(),
			Make:        h.make.Ptr(),
			Model:       h.model.Ptr(),
			Environment: h.environment.Ptr(),
			Appliance:   h.appliance.Ptr(),
			Location:    h.location.Ptr(),
			Rack:        h.rack.Ptr(),
			Rank:        h.rank.Ptr(),
			Slot:        h.slot.Ptr(),
			Type:        ht,
		}.Build(),
	}.Build()

	_, err := h.Metal.UpdateHost(h.Metal.Context(), req)

	return err
}

func (h *Host) remove(glob string) error {
	req := pb.DeleteHostsRequest_builder{
		Zone:    h.zone.Ptr(),
		Cluster: h.cluster.Ptr(),
		Glob:    &glob,
	}.Build()

	_, err := h.Metal.DeleteHosts(h.Metal.Context(), req)

	return err
}

func (a *HostAttr) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name value",
		Short: "Add an " + attribute + " to a " + host,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0], args[1])
		},
	}

	a.zone.Add(cmd.Flags(), host, true)
	a.cluster.Add(cmd.Flags(), host, false)
	a.host.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *HostAttr) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name [value]",
		Short: "Set a " + host + " " + attribute + "'s properties",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			var value string

			if len(args) > 1 {
				value = args[1]
			}

			return a.update(args[0], value)
		},
	}

	a.zone.Add(cmd.Flags(), host, true)
	a.cluster.Add(cmd.Flags(), host, false)
	a.host.Add(cmd.Flags(), attribute, true)
	a.rename.Add(cmd.Flags(), host)

	return cmd
}

func (a *HostAttr) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " [glob]",
		Short: "List one or more " + host + " " + attribute + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	a.zone.Add(cmd.Flags(), host, false)
	a.cluster.Add(cmd.Flags(), host, false)
	a.host.Add(cmd.Flags(), attribute, false)

	return cmd
}

func (a *HostAttr) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " glob",
		Short: "Remove one or more " + host + " " + attribute + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), host, true)
	a.cluster.Add(cmd.Flags(), host, false)
	a.host.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *HostAttr) create(attr string) error {
	req := pb.CreateHostAttrRequest_builder{
		Zone:    a.zone.Ptr(),
		Cluster: a.cluster.Ptr(),
		Host:    a.host.Ptr(),
		Name:    &attr,
	}.Build()

	_, err := a.Metal.CreateHostAttr(a.Metal.Context(), req)

	return err
}

func (a *HostAttr) list(glob string) error {
	type row struct{ Zone, Cluster, Host, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewHostAttrReader(a.zone.Val(), a.cluster.Val(), a.host.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:    resp.GetZone(),
			Cluster: resp.GetCluster(),
			Host:    resp.GetHost(),
			Attr:    resp.GetName(),
			Value:   resp.GetValue(),
		})
	}

	return nil
}

func (a *HostAttr) update(attr, val string) error {
	var value *string

	if val != "" {
		value = &val
	}

	req := pb.UpdateHostAttrRequest_builder{
		Zone:    a.zone.Ptr(),
		Cluster: a.cluster.Ptr(),
		Host:    a.host.Ptr(),
		Name:    &attr,
		Fields: pb.UpdateHostAttrRequest_Fields_builder{
			Name:  a.rename.Ptr(),
			Value: value,
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateHostAttr(a.Metal.Context(), req)

	return err
}

func (a *HostAttr) remove(glob string) error {
	req := pb.DeleteHostAttrsRequest_builder{
		Zone:    a.zone.Ptr(),
		Cluster: a.cluster.Ptr(),
		Host:    a.host.Ptr(),
		Glob:    &glob,
	}.Build()

	_, err := a.Metal.DeleteHostAttrs(a.Metal.Context(), req)

	return err
}
