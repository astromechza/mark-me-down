package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const usageString = `mark-me-down is a simple binary for rendering Github Flavoured Markdown content.

Run it with a single argument (a path to a file) that will be formatted into html
and served whenever a request hits the local server.

The --listen-port field is provided in order to specify a particular port.

`

func mainInner() error {
	sourceFlag := flag.String("source", "", "Source to load from (file path, -, or url")
	sinkFlag := flag.String("sink", "", "Destination to write to (iface:port, file path, -")

	// set a more verbose usage message.
	flag.Usage = func() {
		os.Stderr.WriteString(usageString)
		flag.PrintDefaults()
	}
	// parse them
	flag.Parse()

	// expect a single argument
	if flag.NArg() != 0 {
		return errors.New("no positional arguments allowed")
	}

	if *sourceFlag == "" {
		return errors.New("-source must be provided (see --help)")
	}
	if *sinkFlag == "" {
		return errors.New("-sink must be provided (see --help)")
	}

	var loader SourceLoader
	switch {
	case *sourceFlag == "-":
		content, _ := ioutil.ReadAll(os.Stdin)
		x := OnceOffLoader(content)
		loader = &x
	case strings.HasPrefix(*sourceFlag, "http://") || strings.HasPrefix(*sourceFlag, "https://"):
		x := URLLoader(*sourceFlag)
		loader = &x
	case strings.HasSuffix(*sourceFlag, "/"):
		x := DirectoryLoader(*sourceFlag)
		loader = &x
	default:
		x := FileLoader(*sourceFlag)
		loader = &x
	}

	var sink Sink
	switch {
	case *sinkFlag == "-":
		x := StdoutSink{}
		sink = &x
	case regexp.MustCompile("^[^:]*:\\d+$").MatchString(*sinkFlag):
		x := HttpSink{Address: *sinkFlag}
		sink = &x
	default:
		x := PathSink(*sinkFlag)
		sink = &x
	}
	return sink.Do(loader)
}

func main() {
	err := mainInner()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
