package commands

import (
	"github.com/spf13/cobra"

	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/stack/internal/flags/set"
	"endobit.io/table"
)

type Rack struct {
	*Root
}

func NewRack(r *Root) *Rack {
	return &Rack{Root: r}
}

type RackAttr struct {
	*Rack
	rack set.Rack
}

func NewRackAttr(a *Rack) *RackAttr {
	return &RackAttr{Rack: a}
}

func (a *Rack) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   rack + " name",
		Short: "Add a " + rack + " to a zone",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), rack, true)

	cmd.AddCommand(NewRackAttr(a).Add())

	return cmd
}

func (a *Rack) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   rack + " name",
		Short: "Set a " + rack + "'s properties",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.update(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), rack, true)
	a.rename.Add(cmd.Flags(), rack)

	cmd.AddCommand(NewRackAttr(a).Set())

	return cmd
}

func (a *Rack) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   rack + " [glob]",
		Short: "List one or more " + rack + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	a.zone.Add(cmd.Flags(), rack, false)

	cmd.AddCommand(NewRackAttr(a).List())

	return cmd
}

func (a *Rack) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   rack + " glob",
		Short: "Remove one or more " + rack + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), rack, true)

	cmd.AddCommand(NewRackAttr(a).Remove())

	return cmd
}

func (a *Rack) create(rack string) error {
	req := pb.CreateRackRequest_builder{
		Zone: a.zone.Ptr(),
		Name: &rack,
	}.Build()

	_, err := a.Metal.CreateRack(a.Metal.Context(), req)

	return err
}

func (a *Rack) list(glob string) error {
	type row struct{ Zone, Rack string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewRackReader(a.zone.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone: resp.GetZone(),
			Rack: resp.GetName(),
		})
	}

	return nil
}

func (a *Rack) update(rack string) error {
	req := pb.UpdateRackRequest_builder{
		Zone: a.zone.Ptr(),
		Name: &rack,
		Fields: pb.UpdateRackRequest_Fields_builder{
			Name: a.rename.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateRack(a.Metal.Context(), req)

	return err
}

func (a *Rack) remove(glob string) error {
	req := pb.DeleteRacksRequest_builder{
		Zone: a.zone.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := a.Metal.DeleteRacks(a.Metal.Context(), req)

	return err
}

func (a *RackAttr) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name value",
		Short: "Add an " + attribute + " to a " + rack,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0], args[1])
		},
	}

	a.zone.Add(cmd.Flags(), rack, true)
	a.rack.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *RackAttr) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name [value]",
		Short: "Set a " + rack + " " + attribute + "'s properties",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			var value string

			if len(args) > 1 {
				value = args[1]
			}

			return a.update(args[0], value)
		},
	}

	a.zone.Add(cmd.Flags(), rack, true)
	a.rack.Add(cmd.Flags(), attribute, true)
	a.rename.Add(cmd.Flags(), rack)

	return cmd
}

func (a *RackAttr) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " [glob]",
		Short: "List one or more " + rack + " " + attribute + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	a.zone.Add(cmd.Flags(), rack, false)
	a.rack.Add(cmd.Flags(), attribute, false)

	return cmd
}

func (a *RackAttr) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " glob",
		Short: "Remove one or more " + rack + " " + attribute + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), rack, true)
	a.rack.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *RackAttr) create(attr string) error {
	req := pb.CreateRackAttrRequest_builder{
		Zone: a.zone.Ptr(),
		Rack: a.rack.Ptr(),
		Name: &attr,
	}.Build()

	_, err := a.Metal.CreateRackAttr(a.Metal.Context(), req)

	return err
}

func (a *RackAttr) list(glob string) error {
	type row struct{ Zone, Rack, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewRackAttrReader(a.zone.Val(), a.rack.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:  resp.GetZone(),
			Rack:  resp.GetRack(),
			Attr:  resp.GetName(),
			Value: resp.GetValue(),
		})
	}

	return nil
}

func (a *RackAttr) update(attr, val string) error {
	var value *string

	if val != "" {
		value = &val
	}

	req := pb.UpdateRackAttrRequest_builder{
		Zone: a.zone.Ptr(),
		Rack: a.rack.Ptr(),
		Name: &attr,
		Fields: pb.UpdateRackAttrRequest_Fields_builder{
			Name:  a.rename.Ptr(),
			Value: value,
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateRackAttr(a.Metal.Context(), req)

	return err
}

func (a *RackAttr) remove(glob string) error {
	req := pb.DeleteRackAttrsRequest_builder{
		Zone: a.zone.Ptr(),
		Rack: a.rack.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := a.Metal.DeleteRackAttrs(a.Metal.Context(), req)

	return err
}
