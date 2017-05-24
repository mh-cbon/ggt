package main

import (
	"fmt"

	"github.com/posener/flag"

	"github.com/mh-cbon/ggt/chaner"
	httpProvider "github.com/mh-cbon/ggt/http-provider"
	"github.com/mh-cbon/ggt/mutexer"
	"github.com/mh-cbon/ggt/opts"
	"github.com/mh-cbon/ggt/slicer"
)

var name = "ggt"
var version = "0.0.0"

var subcmds = []string{
	"slicer",
}

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

	}

}

func wrongInput(format string, a ...interface{}) {
	showHelp()
	fmt.Printf(`
    wrong input: %v
    `, fmt.Sprintf(format, a...))
}
func showHelp() {
	fmt.Printf(``)
}
func showVersion() {
	fmt.Printf(`%v - %v
    `, name, version)
}
