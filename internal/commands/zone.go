package commands

import (
	"github.com/spf13/cobra"

	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/stack/internal/flags/set"
	"endobit.io/table"
)

type Zone struct {
	*Root
	timezone set.TimeZone
	template set.Template
}

func NewZone(r *Root) *Zone {
	return &Zone{Root: r}
}

type ZoneAttr struct {
	*Zone
}

func NewZoneAttr(a *Zone) *ZoneAttr {
	return &ZoneAttr{Zone: a}
}

func (z *Zone) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   zone + " name",
		Short: "Add a " + zone,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := z.create(args[0]); err != nil {
				return err
			}

			return z.update(args[0])
		},
	}

	z.timezone.Add(cmd.Flags(), zone)

	return cmd
}

func (z *Zone) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   zone + " name",
		Short: "Set a " + zone + "'s properties",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return z.update(args[0])
		},
	}

	z.timezone.Add(cmd.Flags(), zone)
	z.rename.Add(cmd.Flags(), zone)

	return cmd
}

func (z *Zone) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   zone + " [glob]",
		Short: "List one or more " + zone + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return z.list(glob)
		},
	}

	return cmd
}

func (z *Zone) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   zone + " glob",
		Short: "Remove one or more " + zone + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return z.remove(args[0])
		},
	}

	return cmd
}

func (z *Zone) create(zone string) error {
	var req pb.CreateZoneRequest

	req.SetName(zone)
	_, err := z.Metal.CreateZone(z.Metal.Context(), &req)

	return err
}

func (z *Zone) list(glob string) error {
	type row struct{ Zone, TimeZone string }
	t := table.New()
	defer t.Flush()

	r := z.Metal.NewZoneReader(glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:     resp.GetName(),
			TimeZone: resp.GetTimeZone(),
		})
	}

	return nil
}

func (z *Zone) update(zone string) error {
	req := pb.UpdateZoneRequest_builder{
		Name: &zone,
		Fields: pb.UpdateZoneRequest_Fields_builder{
			Name:     z.rename.Ptr(),
			TimeZone: z.timezone.Ptr(),
		}.Build(),
	}.Build()

	_, err := z.Metal.UpdateZone(z.Metal.Context(), req)

	return err
}

func (z *Zone) remove(glob string) error {
	var req pb.DeleteZonesRequest

	req.SetGlob(glob)
	_, err := z.Metal.DeleteZones(z.Metal.Context(), &req)

	return err
}

func (a *ZoneAttr) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name",
		Short: "Add an " + attribute + " to a " + zone,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), zone, true)

	return cmd
}

func (a *ZoneAttr) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name",
		Short: "Set a " + zone + " " + attribute + "'s properties",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.update(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), zone, true)

	return cmd
}

func (a *ZoneAttr) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " [glob]",
		Short: "List one or more " + zone + " " + attribute + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	a.zone.Add(cmd.Flags(), zone, false)
	a.rename.Add(cmd.Flags(), zone)

	return cmd
}

func (a *ZoneAttr) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " glob",
		Short: "Remove one or more " + zone + " " + attribute + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), zone, true)
	return cmd
}

func (a *ZoneAttr) create(attr string) error {
	req := pb.CreateZoneAttrRequest_builder{
		Zone: a.zone.Ptr(),
		Name: &attr,
	}.Build()

	_, err := a.Metal.CreateZoneAttr(a.Metal.Context(), req)

	return err
}

func (a *ZoneAttr) list(glob string) error {
	type row struct{ Zone, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewZoneAttrReader(a.zone.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:  resp.GetZone(),
			Attr:  resp.GetName(),
			Value: resp.GetValue(),
		})
	}

	return nil
}

func (a *ZoneAttr) update(attr string) error {
	req := pb.UpdateZoneAttrRequest_builder{
		Zone: a.zone.Ptr(),
		Name: &attr,
		Fields: pb.UpdateZoneAttrRequest_Fields_builder{
			Name: a.rename.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateZoneAttr(a.Metal.Context(), req)

	return err
}

func (a *ZoneAttr) remove(glob string) error {
	req := pb.DeleteZoneAttrsRequest_builder{
		Zone: a.zone.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := a.Metal.DeleteZoneAttrs(a.Metal.Context(), req)

	return err
}
