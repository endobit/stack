package commands

import (
	"github.com/spf13/cobra"

	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/stack/internal/flags/set"
	"endobit.io/table"
)

type Environment struct {
	*Root
}

func NewEnvironment(r *Root) *Environment {
	return &Environment{Root: r}
}

type EnvironmentAttr struct {
	*Environment
	environment set.Environment
}

func NewEnvironmentAttr(a *Environment) *EnvironmentAttr {
	return &EnvironmentAttr{Environment: a}
}

func (a *Environment) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   environment + " name",
		Short: "Add an " + environment + " to a zone",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), environment, true)

	cmd.AddCommand(NewEnvironmentAttr(a).Add())

	return cmd
}

func (a *Environment) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   environment + " name",
		Short: "Set an " + environment + "'s properties",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.update(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), environment, true)
	a.rename.Add(cmd.Flags(), environment)

	cmd.AddCommand(NewEnvironmentAttr(a).Set())

	return cmd
}

func (a *Environment) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   environment + " [glob]",
		Short: "List one or more " + environment + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	a.zone.Add(cmd.Flags(), environment, false)

	cmd.AddCommand(NewEnvironmentAttr(a).List())

	return cmd
}

func (a *Environment) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   environment + " glob",
		Short: "Remove one or more " + environment + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), environment, true)

	cmd.AddCommand(NewEnvironmentAttr(a).Remove())

	return cmd
}

func (a *Environment) create(environment string) error {
	req := pb.CreateEnvironmentRequest_builder{
		Zone: a.zone.Ptr(),
		Name: &environment,
	}.Build()

	_, err := a.Metal.CreateEnvironment(a.Metal.Context(), req)

	return err
}

func (a *Environment) list(glob string) error {
	type row struct{ Zone, Environment string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewEnvironmentReader(a.zone.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:        resp.GetZone(),
			Environment: resp.GetName(),
		})
	}

	return nil
}

func (a *Environment) update(environment string) error {
	req := pb.UpdateEnvironmentRequest_builder{
		Zone: a.zone.Ptr(),
		Name: &environment,
		Fields: pb.UpdateEnvironmentRequest_Fields_builder{
			Name: a.rename.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateEnvironment(a.Metal.Context(), req)

	return err
}

func (a *Environment) remove(glob string) error {
	req := pb.DeleteEnvironmentsRequest_builder{
		Zone: a.zone.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := a.Metal.DeleteEnvironments(a.Metal.Context(), req)

	return err
}

func (a *EnvironmentAttr) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name value",
		Short: "Add an " + attribute + " to an " + environment,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0], args[1])
		},
	}

	a.zone.Add(cmd.Flags(), environment, true)
	a.environment.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *EnvironmentAttr) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name [value]",
		Short: "Set an " + environment + " " + attribute + "'s properties",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			var value string

			if len(args) > 1 {
				value = args[1]
			}

			return a.update(args[0], value)
		},
	}

	a.zone.Add(cmd.Flags(), environment, false)
	a.environment.Add(cmd.Flags(), attribute, false)
	a.rename.Add(cmd.Flags(), environment)

	return cmd
}

func (a *EnvironmentAttr) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " [glob]",
		Short: "List one or more " + environment + " " + attribute + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	a.zone.Add(cmd.Flags(), environment, false)
	a.environment.Add(cmd.Flags(), attribute, false)

	return cmd
}

func (a *EnvironmentAttr) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " glob",
		Short: "Remove one or more " + environment + " " + attribute + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), environment, true)
	a.environment.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *EnvironmentAttr) create(attr string) error {
	req := pb.CreateEnvironmentAttrRequest_builder{
		Zone:        a.zone.Ptr(),
		Environment: a.environment.Ptr(),
		Name:        &attr,
	}.Build()

	_, err := a.Metal.CreateEnvironmentAttr(a.Metal.Context(), req)

	return err
}

func (a *EnvironmentAttr) list(glob string) error {
	type row struct{ Zone, Environment, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewEnvironmentAttrReader(a.zone.Val(), a.environment.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:        resp.GetZone(),
			Environment: resp.GetEnvironment(),
			Attr:        resp.GetName(),
			Value:       resp.GetValue(),
		})
	}

	return nil
}

func (a *EnvironmentAttr) update(attr, val string) error {
	var value *string

	if val != "" {
		value = &val
	}

	req := pb.UpdateEnvironmentAttrRequest_builder{
		Zone:        a.zone.Ptr(),
		Environment: a.environment.Ptr(),
		Name:        &attr,
		Fields: pb.UpdateEnvironmentAttrRequest_Fields_builder{
			Name:  a.rename.Ptr(),
			Value: value,
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateEnvironmentAttr(a.Metal.Context(), req)

	return err
}

func (a *EnvironmentAttr) remove(glob string) error {
	req := pb.DeleteEnvironmentAttrsRequest_builder{
		Zone:        a.zone.Ptr(),
		Environment: a.environment.Ptr(),
		Glob:        &glob,
	}.Build()

	_, err := a.Metal.DeleteEnvironmentAttrs(a.Metal.Context(), req)

	return err
}
