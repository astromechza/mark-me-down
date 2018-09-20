package main

import (
	"fmt"
	"github.com/pkg/browser"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

type Sink interface {
	Do(loader SourceLoader) error
}

type StdoutSink struct{}

func (s *StdoutSink) Do(loader SourceLoader) error {
	return loader.Load("", os.Stdout)
}

type HttpSink struct {
	Address string
}

func (s *HttpSink) Do(loader SourceLoader) error {
	// register http function for handling argument /
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		x := strings.TrimLeft(req.URL.Path, "/")
		if err := loader.Load(x, w); err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		}
	})

	// set up server object
	server := &http.Server{Addr: s.Address, Handler: nil}
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}

	tp := ln.(*net.TCPListener)
	fmt.Printf("Listening on %s...\n", tp.Addr())

	// begin serving
	go func() {
		err := server.Serve(ln)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1)
		}
	}()

	// open the server in the available browser
	fmt.Println("Attempting to open a browser window to the address..")
	browser.OpenURL(fmt.Sprintf("http://%s", tp.Addr()))

	// instead of sitting in a for loop or something, we wait for sigint
	signalChannel := make(chan os.Signal, 1)
	// notify that we are going to handle interrupts
	signal.Notify(signalChannel, os.Interrupt)
	for sig := range signalChannel {
		fmt.Printf("Received %v signal. Stopping.\n", sig)
		return nil
	}
	return nil
}

type PathSink string

func (s *PathSink) Do(loader SourceLoader) error {
	switch v := loader.(type) {
	case *DirectoryLoader:
		if !strings.HasSuffix(string(*s), "/") {
			return fmt.Errorf("cannot use directory source with single file sink")
		}
		st, err := os.Stat(string(*s))
		if err != nil {
			return fmt.Errorf("cannot use directory sink when path does not exist")
		}
		if !st.IsDir() {
			return fmt.Errorf("sink path exists but is not a directory")
		}
		return filepath.Walk(string(*v), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fp := filepath.Join(string(*v), path)
			if info.IsDir() {
				if _, err := os.Stat(fp); err != nil {
					if os.IsNotExist(err) {
						return os.Mkdir(fp, 0644)
					}
					return err
				}
			} else {
				if !strings.HasSuffix(fp, ".md") || !strings.HasSuffix(fp, ".MD") {
					return nil
				}

				f, err := os.Create(fp)
				if err != nil {
					return err
				}
				err = v.Load(path, f)
				f.Close()
				return err
			}
			return nil
		})
	default:
		f, err := os.Create(string(*s))
		if err != nil {
			return err
		}
		defer f.Close()
		return v.Load("", f)
	}
}
