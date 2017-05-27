// Package opts of the cli.
package opts

import "github.com/posener/flag"

// Cli struct for options
type Cli struct {
	Help     bool
	Version  bool
	Verbose  bool
	Contract bool
	OutPkg   string
	Mode     string
	Args     []string
}

// Bind options on flag
func (c *Cli) Bind() {

	flag.BoolVar(&c.Help, "help", false, "show help")
	flag.BoolVar(&c.Version, "version", false, "show version")
	flag.BoolVar(&c.Verbose, "vv", false, "more verbose")
	flag.BoolVar(&c.Contract, "c", false, "with contract if the generator supports it")
	flag.StringVar(&c.OutPkg, "p", "", "Package name of the new code.")
	flag.StringVar(&c.Mode, "mode", "", "Generator mode when suitable (rpc|route)")

}
