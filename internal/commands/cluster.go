package commands

import (
	"errors"

	"github.com/spf13/cobra"

	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/table"
)

type Cluster struct {
	*Root
}

func NewCluster(r *Root) *Cluster {
	return &Cluster{Root: r}
}

type ClusterAttr struct {
	*Cluster
}

func NewClusterAttr(c *Cluster) *ClusterAttr {
	return &ClusterAttr{Cluster: c}
}

func (a *Cluster) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   cluster + " name",
		Short: "Add a " + cluster + " to a zone",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), cluster, true)

	cmd.AddCommand(NewClusterAttr(a).Add())

	return cmd
}

func (a *Cluster) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   cluster + " name",
		Short: "Set a " + cluster + "'s properties",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.update(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), cluster, true)
	a.rename.Add(cmd.Flags(), cluster)

	cmd.AddCommand(NewClusterAttr(a).Set())

	return cmd
}

func (a *Cluster) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   cluster + " [glob]",
		Short: "List one or more " + cluster + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
				if !a.zone.IsSet() {
					return errors.New("zone must be specified when using glob pattern")
				}
			}
			return a.list(glob)
		},
	}

	a.zone.Add(cmd.Flags(), cluster, false)

	cmd.AddCommand(NewClusterAttr(a).List())

	return cmd
}

func (a *Cluster) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   cluster + " glob",
		Short: "Remove one or more " + cluster + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), cluster, true)

	cmd.AddCommand(NewClusterAttr(a).Remove())

	return cmd
}

func (a *Cluster) create(cluster string) error {
	req := pb.CreateClusterRequest_builder{
		Zone: a.zone.Ptr(),
		Name: &cluster,
	}.Build()

	_, err := a.Metal.CreateCluster(a.Metal.Context(), req)

	return err
}

func (a *Cluster) list(glob string) error {
	type row struct{ Zone, Cluster string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewClusterReader(a.zone.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:    resp.GetZone(),
			Cluster: resp.GetName(),
		})
	}

	return nil
}

func (a *Cluster) update(cluster string) error {
	req := pb.UpdateClusterRequest_builder{
		Zone: a.zone.Ptr(),
		Name: &cluster,
		Fields: pb.UpdateClusterRequest_Fields_builder{
			Name: a.rename.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateCluster(a.Metal.Context(), req)

	return err
}

func (a *Cluster) remove(glob string) error {
	req := pb.DeleteClustersRequest_builder{
		Zone: a.zone.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := a.Metal.DeleteClusters(a.Metal.Context(), req)

	return err
}

func (a *ClusterAttr) Add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name value",
		Short: "Add an " + attribute + " to a " + cluster,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := a.create(args[0]); err != nil {
				return err
			}

			return a.update(args[0], args[1])
		},
	}

	a.zone.Add(cmd.Flags(), cluster, true)
	a.cluster.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *ClusterAttr) Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " name [value]",
		Short: "Set a " + cluster + " " + attribute + "'s properties",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			var value string

			if len(args) > 1 {
				value = args[1]
			}

			return a.update(args[0], value)
		},
	}

	a.zone.Add(cmd.Flags(), cluster, true)
	a.cluster.Add(cmd.Flags(), attribute, true)
	a.rename.Add(cmd.Flags(), cluster)

	return cmd
}

func (a *ClusterAttr) List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " [glob]",
		Short: "List one or more " + cluster + " " + attribute + "s",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) > 0 {
				glob = args[0]
			}
			return a.list(glob)
		},
	}

	a.zone.Add(cmd.Flags(), cluster, false)
	a.cluster.Add(cmd.Flags(), attribute, false)

	return cmd
}

func (a *ClusterAttr) Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   attribute + " glob",
		Short: "Remove one or more " + cluster + " " + attribute + "s",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.remove(args[0])
		},
	}

	a.zone.Add(cmd.Flags(), cluster, true)
	a.cluster.Add(cmd.Flags(), attribute, true)

	return cmd
}

func (a *ClusterAttr) create(attr string) error {
	req := pb.CreateClusterAttrRequest_builder{
		Zone:    a.zone.Ptr(),
		Cluster: a.cluster.Ptr(),
		Name:    &attr,
	}.Build()

	_, err := a.Metal.CreateClusterAttr(a.Metal.Context(), req)

	return err
}

func (a *ClusterAttr) list(glob string) error {
	type row struct{ Zone, Cluster, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Metal.NewClusterAttrReader(a.zone.Val(), a.cluster.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:    resp.GetZone(),
			Cluster: resp.GetCluster(),
			Attr:    resp.GetName(),
			Value:   resp.GetValue(),
		})
	}

	return nil
}

func (a *ClusterAttr) update(attr, val string) error {
	var value *string

	if val != "" {
		value = &val
	}

	req := pb.UpdateClusterAttrRequest_builder{
		Zone:    a.zone.Ptr(),
		Cluster: a.cluster.Ptr(),
		Name:    &attr,
		Fields: pb.UpdateClusterAttrRequest_Fields_builder{
			Name:  a.rename.Ptr(),
			Value: value,
		}.Build(),
	}.Build()

	_, err := a.Metal.UpdateClusterAttr(a.Metal.Context(), req)

	return err
}

func (a *ClusterAttr) remove(glob string) error {
	req := pb.DeleteClusterAttrsRequest_builder{
		Zone:    a.zone.Ptr(),
		Cluster: a.cluster.Ptr(),
		Glob:    &glob,
	}.Build()

	_, err := a.Metal.DeleteClusterAttrs(a.Metal.Context(), req)

	return err
}
