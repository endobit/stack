package stack

//go:generate go run "github.com/dmarkham/enumer" -type HostType -linecomment -text

// HostType is the type of host.
type HostType int

// Host types.
const (
	MetalHostType     HostType = iota // metal
	VirtualHostType                   // virtual
	ContainerHostType                 // container
)
