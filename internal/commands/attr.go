package commands

import (
	"github.com/spf13/cobra"

	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/table"
)

type Attr struct {
	*Root
}

func NewAttr(r *Root) *Attr {
	return &Attr{Root: r}
}

func (a *Attr) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name",
		Short: "Add a global " + attribute,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0])
		},
	}

	return cmd
}

func (a *Attr) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name",
		Short: "Set a global " + attribute + "'s properties",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.update(args[0])
		},
	}

	a.rename.Add(cmd.Flags(), attribute)

	return cmd
}

func (a *Attr) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " [glob]",
		Short: "List one or more global " + attribute + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	return cmd
}

func (a *Attr) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " glob",
		Short: "Remove one or more global " + attribute + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	return cmd
}

func (a *Attr) create(attr string) error {
	req := pb.CreateGlobalAttrRequest_builder{
		Name: &attr,
	}.Build()

	_, err := a.Metal.CreateGlobalAttr(a.Metal.Context(), req)

	return err
}

func (a *Attr) list(glob string) error {
	type row struct{ Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewGlobalAttrReader(glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Attr:  resp.GetName(),
			Value: resp.GetValue(),
		})
	}

	return nil
}

func (a *Attr) update(attr string) error {
	req := pb.UpdateGlobalAttrRequest_builder{
		Name: &attr,
		Fields: pb.UpdateGlobalAttrRequest_Fields_builder{
			Name: a.rename.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateGlobalAttr(a.Metal.Context(), req)

	return err
}

func (a *Attr) remove(glob string) error {
	req := pb.DeleteGlobalAttrsRequest_builder{
		Glob: &glob,
	}.Build()

	_, err := a.Metal.DeleteGlobalAttrs(a.Metal.Context(), req)

	return err
}
