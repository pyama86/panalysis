package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pyama86/panalysis/panalysis"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

var (
	version   string
	revision  string
	goversion string
	builddate string
	builduser string
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

func printVersion() {
	fmt.Printf("panalysis version: %s (%s)\n", version, revision)
	fmt.Printf("build at %s (with %s) by %s\n", builddate, goversion, builduser)
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		jsonp     bool
		config    bool
		arguments []string
		version   bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet("panalysis", flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&jsonp, "json", false, "Json")
	flags.BoolVar(&jsonp, "j", false, "Json(Short)")
	flags.BoolVar(&config, "config", false, "Config")
	flags.BoolVar(&config, "c", false, "Config(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		printVersion()
		return ExitCodeOK
	}
	var fd io.Reader
	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}
	if len(os.Args) < 2 {
		fd = os.Stdin
	} else {
		fp, err := os.Open(arguments[0])
		if err != nil {
			log.Fatal(err)
		}
		fd = fp
		defer fp.Close()
	}
	var pn panalysis.Parser
	if !jsonp {
		pn = panalysis.NewConfigParser(fd)
	} else {
		pn = panalysis.NewJSONParser(fd)
	}

	r, err := pn.Parse()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(os.Stdout, r)
	return ExitCodeOK
}
