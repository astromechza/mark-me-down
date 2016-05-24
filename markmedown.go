package main

import (
    "os"
    "flag"
    "fmt"
    "net/http"
    "io/ioutil"
    "os/signal"

    "github.com/pkg/browser"
    "github.com/russross/blackfriday"
)

const usageString =
`mark-me-down is a simple binary for rendering Github Flavoured Markdown content.

Run it with a single argument (a path to a file) that will be formatted into html
and served whenever a request hits the local server.

The --listen-port field is provided in order to specify a particular port.

`

// address to listen on
const listenAddress = "localhost"

// html template for wrapping the markdown html into proper html
const htmlTemplate = "<html><head><title>%s</title><style>%s</style></head><body class=\"markdown-body\">%s</body></html>"

// return a function that formats the content of the given file path on request
func buildMarkdownFileServer(filepath string) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        content, err := ioutil.ReadFile(filepath)
        if err != nil {
            fmt.Fprintf(w, err.Error())
        } else {
            formatted := string(blackfriday.MarkdownCommon(content))
            fmt.Fprintf(w, fmt.Sprintf(htmlTemplate, filepath, GFMCSS, formatted))
        }
    }
}

func main() {
    // first set up config flag options
    listenPortFlag := flag.Int("listen-port", 80, "Server the markdown on this port")

    // set a more verbose usage message.
    flag.Usage = func() {
        os.Stderr.WriteString(usageString)
        flag.PrintDefaults()
    }
    // parse them
    flag.Parse()

    // expect a single argument
    if flag.NArg() != 1 {
        os.Stderr.WriteString("A single input-file is required as argument 1. Use --help to see the usage.\n")
        os.Exit(1)
    }

    // a listen port
    if *listenPortFlag < 1 {
        os.Stderr.WriteString("listen-port must be > 0. Use --help to see the usage.\n")
        os.Exit(1)
    }

    // register http function for handling argument /
    http.HandleFunc("/", buildMarkdownFileServer(flag.Args()[0]))

    // begin serving
    go func() {
        err := http.ListenAndServe(fmt.Sprintf("%s:%d", listenAddress, *listenPortFlag), nil)
        if err != nil {
            os.Stderr.WriteString(err.Error() + "\n")
            os.Exit(1)
        }
    }()

    fmt.Printf("Listening on %s:%d...\n", listenAddress, *listenPortFlag)

    // open the server in the available browser
    fmt.Println("Attempting to open a browser window to the address..")
    browser.OpenURL(fmt.Sprintf("http://%s:%d", listenAddress, *listenPortFlag))

    // instead of sitting in a for loop or something, we wait for sigint
    signalChannel := make(chan os.Signal, 1)
    // notify that we are going to handle interrupts
    signal.Notify(signalChannel, os.Interrupt)
    for sig := range signalChannel {
        fmt.Printf("Received %v signal. Stopping.\n", sig)
        os.Exit(0)
    }
}
