// ggt's generator toolbox
package main

import (
	"fmt"

	"github.com/posener/flag"

	"github.com/mh-cbon/ggt/chaner"
	httpConsumer "github.com/mh-cbon/ggt/http-consumer"
	httpProvider "github.com/mh-cbon/ggt/http-provider"
	"github.com/mh-cbon/ggt/mutexer"
	"github.com/mh-cbon/ggt/opts"
	"github.com/mh-cbon/ggt/slicer"
)

var name = "ggt"
var version = "0.0.0"

func main() {

	options := &opts.Cli{}
	options.Bind()

	flag.SetInstallFlags("complete", "uncomplete")
	flag.Parse()
	if flag.Complete() {
		return
	}

	if flag.NArg() < 1 {
		if options.Help {
			showVersion()
			showHelp()
			return
		} else if options.Version {
			showVersion()
			return
		} else {
			wrongInput("wrong invokation")
			return
		}
	}

	args := flag.Args()
	cmd := args[0]
	options.Args = args[1:]

	if cmd == "slicer" {
		(slicer.Cmd{}).Run(options)

	} else if cmd == "mutexer" {
		(mutexer.Cmd{}).Run(options)

	} else if cmd == "chaner" {
		(chaner.Cmd{}).Run(options)

	} else if cmd == "http-provider" {
		(httpProvider.Cmd{}).Run(options)

	} else if cmd == "http-consumer" {
		(httpConsumer.Cmd{}).Run(options)

	}

}

func wrongInput(format string, a ...interface{}) {
	showHelp()
	fmt.Printf(`
    wrong input: %v
    `, fmt.Sprintf(format, a...))
}
func showHelp() {
	fmt.Printf(`ggt [options] [generator] [...types]

ggt's generator toolbox

[options]
    -help        Show help
    -version     Show version
    -vv          More verbose
    -mode        Generator mode when suitable (rpc|route).

[generator]

    One of slicer, chaner, mutexer, http-provider.

[...types]
    A list of types such as src:dst.
    A type is defined by its package path and its type name,
    [pkgpath/]name.
    If the Package path is empty, it is set to the package name being generated.
    Name can be a valid type identifier such as TypeName, *TypeName, []TypeName

Example

    ggt -c slicer MySrcType:gen/*NewGenType
    ggt -c slicer myModule/*MySrcType:gen/NewGenType
`)
}
func showVersion() {
	fmt.Printf(`%v - %v
    `, name, version)
}
