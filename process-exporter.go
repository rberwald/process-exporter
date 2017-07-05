package main

import (
//	"bytes"
//	"encoding/json"
//	"flag"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"net/url"
//	"os"
//	"regexp"
//	"strconv"
//	"strings"
//	"time"
//
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

const (
	namespace = "process"
)

var (
)

type Exporter struct {
	process_metrics	map[int]*prometheus.GaugeVec
}

// Main function
func main() {
	const processHelpText = `Processes to (no)watch
	You should provided at least one process to watch.
	The parameter process.watch should be a comma-seperated list of regular expressions of processes to watch
	The parameter process.nowatch is a filter that removes processes from the list provided by process.watch`

	var (
		listenAddress   = flag.String("web.listen-address", ":8980", "Address to listen on for web interface and telemetry.")
		metricsPath     = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		showVersion     = flag.Bool("version", false, "Print version information.")
		processExpr	= flag.String("process.watch", "", processHelpText)
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("process_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting process_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

//	exporter, err := rooomyExporter(*rooomyScrapeURI, rooomy_env_user, rooomy_env_password, rooomy_env_environment, rooomy_env_role)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	prometheus.MustRegister(exporter)
//	prometheus.MustRegister(version.NewCollector("rooomy_exporter"))
//
//	log.Infoln("Listening on", *listenAddress)
//	http.Handle(*metricsPath, prometheus.Handler())
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte(`<html>
//             <head><title>Rooomy Exporter</title></head>
//             <body>
//             <h1>Rooomy Exporter</h1>
//             <p><a href='` + *metricsPath + `'>Metrics</a></p>
//             </body>
//             </html>`))
//	})
//	log.Fatal(http.ListenAndServe(*listenAddress, nil))
//
}
