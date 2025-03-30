package commands

import (
	"github.com/spf13/cobra"

	"endobit.io/metal"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/stack/internal/flags/set"
	"endobit.io/table"
)

type Model struct {
	*Root
	make set.Make
	arch set.Arch
}

func NewModel(r *Root) *Model {
	return &Model{Root: r}
}

type ModelAttr struct {
	*Model
	model set.Model
}

func NewModelAttr(a *Model) *ModelAttr {
	return &ModelAttr{Model: a}
}

func (m *Model) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   model + " make name",
		Short: "Add a " + model,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := m.create(args[0], args[1]); err != nil {
				return err
			}

			return m.update(args[0], args[1])
		},
	}

	m.arch.Add(cmd.Flags(), model)

	cmd.AddCommand(NewModelAttr(m).Add())

	return cmd
}

func (m *Model) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   model + " make name",
		Short: "Set a " + model + "'s properties",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return m.update(args[0], args[1])
		},
	}

	m.arch.Add(cmd.Flags(), model)
	m.rename.Add(cmd.Flags(), model)

	cmd.AddCommand(NewModelAttr(m).Set())

	return cmd
}

func (m *Model) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   model + " [make] [glob]",
		Short: "List one or more " + model + "s",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			var vendor, glob string

			if len(args) > 0 {
				vendor = args[0]
			}

			if len(args) > 1 {
				glob = args[1]
			}

			return m.list(vendor, glob)
		},
	}

	cmd.AddCommand(NewModelAttr(m).List())

	return cmd
}

func (m *Model) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   model + " make glob",
		Short: "Remove one or more " + model + "s",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return m.remove(args[0], args[1])
		},
	}

	cmd.AddCommand(NewModelAttr(m).Remove())

	return cmd
}

func (m *Model) create(vendor, model string) error {
	req := pb.CreateModelRequest_builder{
		Make: &vendor,
		Name: &model,
	}.Build()

	_, err := m.Metal.CreateModel(m.Metal.Context(), req)
	return err
}

func (m *Model) list(vendor, glob string) error {
	type row struct{ Make, Model, Arch string }
	t := table.New()
	defer t.Flush()

	r := m.Metal.NewModelReader(vendor, glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Make:  resp.GetMake(),
			Model: resp.GetName(),
			Arch:  metal.ShortArchitecture[resp.GetArchitecture().String()],
		})

	}

	return nil
}

func (m *Model) update(vendor, model string) error {
	var pbarch *pb.Architecture

	if m.arch.Val() != "" {
		a := pb.Architecture(pb.Architecture_value[m.arch.Val()])
		pbarch = &a
	}

	req := pb.UpdateModelRequest_builder{
		Make: &vendor,
		Name: &model,
		Fields: pb.UpdateModelRequest_Fields_builder{
			Name:         m.rename.Ptr(),
			Architecture: pbarch,
		}.Build(),
	}.Build()

	_, err := m.Metal.UpdateModel(m.Metal.Context(), req)

	return err
}

func (m *Model) remove(vendor, glob string) error {
	req := pb.DeleteModelsRequest_builder{
		Make: &vendor,
		Glob: &glob,
	}.Build()

	_, err := m.Metal.DeleteModels(m.Metal.Context(), req)
	return err
}

func (a *ModelAttr) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name value",
		Short: "Add an " + attribute + " to a " + model,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0], args[1])
		},
	}

	a.make.Add(cmd.Flags(), model, true)
	a.model.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *ModelAttr) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name",
		Short: "Set a " + model + " " + attribute + "'s properties",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			var value string

			if len(args) > 1 {
				value = args[1]
			}

			return a.update(args[0], value)
		},
	}

	a.make.Add(cmd.Flags(), model, false)
	a.model.Add(cmd.Flags(), attribute, false)
	a.rename.Add(cmd.Flags(), model)

	return cmd
}

func (a *ModelAttr) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " [glob]",
		Short: "List one or more " + model + " " + attribute + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	a.make.Add(cmd.Flags(), model, true)
	a.model.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *ModelAttr) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " glob",
		Short: "Remove one or more " + model + " " + attribute + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.make.Add(cmd.Flags(), model, true)
	a.model.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *ModelAttr) create(attr string) error {
	req := pb.CreateModelAttrRequest_builder{
		Model: a.model.Ptr(),
		Name:  &attr,
	}.Build()

	_, err := a.Metal.CreateModelAttr(a.Metal.Context(), req)

	return err
}

func (a *ModelAttr) list(glob string) error {
	type row struct{ Make, Model, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewModelAttrReader(a.model.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Make:  resp.GetMake(),
			Model: resp.GetModel(),
			Attr:  resp.GetName(),
			Value: resp.GetValue(),
		})
	}

	return nil
}

func (a *ModelAttr) update(attr, val string) error {
	var value *string

	if val != "" {
		value = &val
	}

	req := pb.UpdateModelAttrRequest_builder{
		Model: a.model.Ptr(),
		Name:  &attr,
		Fields: pb.UpdateModelAttrRequest_Fields_builder{
			Name:  a.rename.Ptr(),
			Value: value,
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateModelAttr(a.Metal.Context(), req)

	return err
}

func (a *ModelAttr) remove(glob string) error {
	req := pb.DeleteModelAttrsRequest_builder{
		Model: a.model.Ptr(),
		Glob:  &glob,
	}.Build()

	_, err := a.Metal.DeleteModelAttrs(a.Metal.Context(), req)

	return err
}
