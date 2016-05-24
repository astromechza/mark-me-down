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
`<mark-me-down usage message>
`

const htmlTemplate = "<html><head><title>%s</title><style>%s</style></head><body class=\"markdown-body\">%s</body></html>"

func formatHTML(markdownFile string, markdownHTML string) string {
    return fmt.Sprintf(htmlTemplate, markdownFile, GFMCSS, markdownHTML)
}

func buildMarkdownFileServer(filepath string) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        content, err := ioutil.ReadFile(filepath)
        if err != nil {
            fmt.Fprintf(w, err.Error())
        } else {
            formatted := string(blackfriday.MarkdownCommon(content))
            fmt.Fprintf(w, formatHTML(filepath, formatted))
        }
    }
}

func mainInner() error {
    // first set up config flag options
    inputFileFlag := flag.String("input-file", "", "The input file to watch and process")
    listenPortFlag := flag.Int("listen-port", 80, "Server the markdown on this port")

    // set a more verbose usage message.
    flag.Usage = func() {
        os.Stderr.WriteString(usageString)
        flag.PrintDefaults()
    }
    // parse them
    flag.Parse()

    if *inputFileFlag == "" {
        return fmt.Errorf("input-file is required")
    }

    http.HandleFunc("/", buildMarkdownFileServer(*inputFileFlag))

    if *listenPortFlag < 1 {
        return fmt.Errorf("listen-port must be > 0")
    }

    go func() {
        err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", *listenPortFlag), nil)
        if err != nil {
            os.Stderr.WriteString(err.Error() + "\n")
            os.Exit(1)
        }
    }()

    browser.OpenURL(fmt.Sprintf("http://127.0.0.1:%d", *listenPortFlag))

    // instead of sitting in a for loop or something, we wait for sigint
    signalChannel := make(chan os.Signal, 1)
    // notify that we are going to handle interrupts
    signal.Notify(signalChannel, os.Interrupt)
    for sig := range signalChannel {
        fmt.Printf("Received %v signal. Stopping.\n", sig)
        os.Exit(0)
    }
    return nil
}

func main() {
    if err := mainInner(); err != nil {
        os.Stderr.WriteString(err.Error() + "\n")
        os.Exit(1)
    }
}
