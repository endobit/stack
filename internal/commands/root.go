package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"

	"endobit.io/metal"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/mops"
	"endobit.io/stack/internal/flags/set"
)

type Root struct {
	Metal   *metal.Client
	Ops     *mops.Client
	zone    set.Zone
	cluster set.Cluster
	host    set.Host
	json    set.JSON
	rename  set.Rename
}

func (r *Root) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	attr := NewAttr(r)
	appliance := NewAppliance(r)
	environment := NewEnvironment(r)
	cluster := NewCluster(r)
	host := NewHost(r)
	rack := NewRack(r)
	model := NewModel(r)
	zone := NewZone(r)

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:     "add",
			Aliases: []string{"create"},
			Short:   "Add objects",
		}

		cmd.AddCommand(
			attr.Add(),
			appliance.Add(),
			cluster.Add(),
			environment.Add(),
			host.Add(),
			model.Add(),
			rack.Add(),
			zone.Add())

	case Dump:
		cmd = cobra.Command{
			Use:   "dump",
			Short: "Dump stack schema",
			Args:  cobra.NoArgs,
			RunE: func(_ *cobra.Command, _ []string) error {
				return r.dump()
			},
		}

		r.json.Add(cmd.Flags(), "schema")
		r.zone.Add(cmd.Flags(), "schema", false)
		r.cluster.Add(cmd.Flags(), "schema", false)
		r.host.Add(cmd.Flags(), "schema", false)

	case Set:
		cmd = cobra.Command{
			Use:     "set",
			Aliases: []string{"update"},
			Short:   "Set object properties",
		}

		cmd.AddCommand(
			attr.Set(),
			appliance.Set(),
			cluster.Set(),
			environment.Set(),
			host.Set(),
			model.Set(),
			rack.Set(),
			zone.Set())

	case Unset:
		cmd = cobra.Command{
			Use:   "unset",
			Short: "Unset object properties",
		}

		cmd.AddCommand(
			host.Unset(),
		)

	case List:
		cmd = cobra.Command{
			Use:     "list",
			Aliases: []string{"ls"},
			Short:   "List objects",
			Long:    "List is for humans.",
		}

		cmd.AddCommand(
			attr.List(),
			appliance.List(),
			cluster.List(),
			environment.List(),
			host.List(),
			model.List(),
			rack.List(),
			zone.List())

	case Load:
		cmd = cobra.Command{
			Use:     "load filename",
			Aliases: []string{"ld"},
			Args:    cobra.ExactArgs(1),
			Short:   "Load objects",
			RunE: func(_ *cobra.Command, args []string) error {
				return r.load(args[0])
			},
		}

	case Report:
		cmd = cobra.Command{
			Use:   "report name",
			Short: "Report from template",
			Long:  "Renders a report from the named template.",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return r.report(args[0])
			},
		}

		r.zone.Add(cmd.Flags(), "report", false)
		r.cluster.Add(cmd.Flags(), "report", false)
		r.host.Add(cmd.Flags(), "report", false)

	case Remove:
		cmd = cobra.Command{
			Use:     "remove",
			Aliases: []string{"del", "rm"},
			Short:   "Remove objects",
		}

		cmd.AddCommand(
			attr.Remove(),
			appliance.Remove(),
			cluster.Remove(),
			environment.Remove(),
			host.Remove(),
			model.Remove(),
			rack.Remove(),
			zone.Remove())
	}

	return &cmd
}

func (r *Root) dump() error {
	req := pb.ReadSchemaRequest_builder{
		Zone:    r.zone.Ptr(),
		Cluster: r.cluster.Ptr(),
		Host:    r.host.Ptr(),
	}.Build()

	resp, err := r.Metal.ReadSchema(r.Metal.Context(), req)
	if err != nil {
		return err
	}

	doc := resp.GetSchema()

	if !r.json.Val() { // parse json as yaml and re-marshal
		var obj map[string]any

		b, err := protojson.MarshalOptions{
			UseProtoNames: true,
		}.Marshal(doc)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(b, &obj); err != nil {
			return err
		}

		return yaml.NewEncoder(os.Stdout, yaml.IndentSequence(true)).Encode(obj)
	}

	b, err := protojson.MarshalOptions{
		Multiline:     true,
		UseProtoNames: true,
	}.Marshal(doc)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func (r *Root) load(filename string) error {
	fin, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fin.Close()

	data, err := io.ReadAll(fin)
	if err != nil {
		return err
	}

	var doc pb.Schema

	switch filepath.Ext(filename) {
	case ".json":
		if err := protojson.Unmarshal(data, &doc); err != nil {
			return err
		}
	case ".yaml", ".yml":
		var jsonMap map[string]any

		if err := yaml.Unmarshal(data, &jsonMap); err != nil {
			return fmt.Errorf("failed to unmarshal yaml: %w", err)
		}

		jsonData, err := json.Marshal(jsonMap)
		if err != nil {
			return fmt.Errorf("failed to marshal json: %w", err)
		}

		if err := protojson.Unmarshal(jsonData, &doc); err != nil {
			return fmt.Errorf("failed to unmarshal json: %w", err)
		}

	default:
		return errors.New("unknown file type")
	}

	req := pb.CreateSchemaRequest_builder{
		Schema: &doc,
	}.Build()

	_, err = r.Metal.CreateSchema(r.Metal.Context(), req)

	return err
}

func (r *Root) report(template string) error {
	scope := mops.ReportScope{
		Zone:    r.zone.Val(),
		Cluster: r.cluster.Val(),
		Host:    r.host.Val(),
	}

	var resp mops.GetReportResponse

	if err := r.Ops.Get(context.Background(), "report/"+template+scope.Query(), &resp); err != nil {
		return fmt.Errorf("failed to get report template %q: %w", template, err)
	}

	fmt.Println(resp.Report)

	return nil
}
