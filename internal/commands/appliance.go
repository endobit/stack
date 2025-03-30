package commands

import (
	"github.com/spf13/cobra"

	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/stack/internal/flags/set"
	"endobit.io/table"
)

type Appliance struct {
	*Root
}

func NewAppliance(r *Root) *Appliance {
	return &Appliance{Root: r}
}

type ApplianceAttr struct {
	*Appliance
	appliance set.Appliance
}

func NewApplianceAttr(a *Appliance) *ApplianceAttr {
	return &ApplianceAttr{Appliance: a}
}

func (a *Appliance) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   appliance + " name",
		Short: "Add an " + appliance + " to a zone",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), appliance, true)

	cmd.AddCommand(NewApplianceAttr(a).Add())

	return cmd
}

func (a *Appliance) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   appliance + " name",
		Short: "Set an " + appliance + "'s properties",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.update(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), appliance, true)
	a.rename.Add(cmd.Flags(), appliance)

	cmd.AddCommand(NewApplianceAttr(a).Set())

	return cmd
}

func (a *Appliance) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   appliance + " [glob]",
		Short: "List one or more " + appliance + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	a.zone.Add(cmd.Flags(), appliance, false)

	cmd.AddCommand(NewApplianceAttr(a).List())

	return cmd
}

func (a *Appliance) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   appliance + " glob",
		Short: "Remove one or more " + appliance + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), appliance, true)

	cmd.AddCommand(NewApplianceAttr(a).Remove())

	return cmd
}

func (a *Appliance) create(appliance string) error {
	req := pb.CreateApplianceRequest_builder{
		Zone: a.zone.Ptr(),
		Name: &appliance,
	}.Build()

	_, err := a.Metal.CreateAppliance(a.Metal.Context(), req)

	return err
}

func (a *Appliance) list(glob string) error {
	type row struct{ Zone, Appliance string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewApplianceReader(a.zone.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:      resp.GetZone(),
			Appliance: resp.GetName(),
		})
	}

	return nil
}

func (a *Appliance) update(appliance string) error {
	req := pb.UpdateApplianceRequest_builder{
		Zone: a.zone.Ptr(),
		Name: &appliance,
		Fields: pb.UpdateApplianceRequest_Fields_builder{
			Name: a.rename.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateAppliance(a.Metal.Context(), req)

	return err
}

func (a *Appliance) remove(glob string) error {
	req := pb.DeleteAppliancesRequest_builder{
		Zone: a.zone.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := a.Metal.DeleteAppliances(a.Metal.Context(), req)

	return err
}

func (a *ApplianceAttr) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name value",
		Short: "Add an " + attribute + " to an " + appliance,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0], args[1])
		},
	}

	a.zone.Add(cmd.Flags(), appliance, true)
	a.appliance.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *ApplianceAttr) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name [value]",
		Short: "Set an " + appliance + " " + attribute + "'s properties",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			var value string

			if len(args) > 1 {
				value = args[1]
			}

			return a.update(args[0], value)
		},
	}

	a.zone.Add(cmd.Flags(), appliance, true)
	a.appliance.Add(cmd.Flags(), attribute, true)
	a.rename.Add(cmd.Flags(), attribute)

	return cmd
}

func (a *ApplianceAttr) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " [glob]",
		Short: "List one or more " + appliance + " " + attribute + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	a.zone.Add(cmd.Flags(), appliance, false)
	a.appliance.Add(cmd.Flags(), attribute, false)

	return cmd
}

func (a *ApplianceAttr) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " glob",
		Short: "Remove one or more " + appliance + " " + attribute + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), appliance, true)
	a.appliance.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *ApplianceAttr) create(attr string) error {
	req := pb.CreateApplianceAttrRequest_builder{
		Zone:      a.zone.Ptr(),
		Appliance: a.appliance.Ptr(),
		Name:      &attr,
	}.Build()

	_, err := a.Metal.CreateApplianceAttr(a.Metal.Context(), req)

	return err
}

func (a *ApplianceAttr) list(glob string) error {
	type row struct{ Zone, Appliance, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewApplianceAttrReader(a.zone.Val(), a.appliance.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:      resp.GetZone(),
			Appliance: resp.GetAppliance(),
			Attr:      resp.GetName(),
			Value:     resp.GetValue(),
		})
	}

	return nil
}

func (a *ApplianceAttr) update(attr, val string) error {
	var value *string

	if val != "" {
		value = &val
	}

	req := pb.UpdateApplianceAttrRequest_builder{
		Zone:      a.zone.Ptr(),
		Appliance: a.appliance.Ptr(),
		Name:      &attr,
		Fields: pb.UpdateApplianceAttrRequest_Fields_builder{
			Name:  a.rename.Ptr(),
			Value: value,
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateApplianceAttr(a.Metal.Context(), req)

	return err
}

func (a *ApplianceAttr) remove(glob string) error {
	req := pb.DeleteApplianceAttrsRequest_builder{
		Zone:      a.zone.Ptr(),
		Appliance: a.appliance.Ptr(),
		Glob:      &glob,
	}.Build()

	_, err := a.Metal.DeleteApplianceAttrs(a.Metal.Context(), req)

	return err
}
