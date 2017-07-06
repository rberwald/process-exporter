package main

import (
	//	"bytes"
	//	"encoding/json"
	"flag"
	"fmt"
	//	"io/ioutil"
	"net/http"
	//	"net/url"
	"os"
	//	"regexp"
	"strconv"
	"strings"
	//	"time"
	//
	"github.com/shirou/gopsutil/process"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

const (
	namespace = "indivdual"
)

var (
	to_watch = []string{}
	to_skip  = []string{}
)

type Exporter struct {
	up             prometheus.Gauge
	totalScrapes   prometheus.Counter
	processMetrics map[string]*prometheus.GaugeVec
}

// NewExporter returns an initialized Exporter.
func NewExporter(processMetrics map[string]*prometheus.GaugeVec) (*Exporter, error) {

	return &Exporter{
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the last scrape of processes successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_total_scrapes",
			Help:      "Current total process scrapes.",
		}),
		processMetrics: processMetrics,
	}, nil
}

func (e *Exporter) scrape() {
	e.totalScrapes.Inc()
	e.up.Set(0)

	procs, err := process.Pids()
	if err != nil {
		log.Errorf("Couldn't retrieve processes: %v", err)
		return
	}

	for _, pid := range procs {
		proc, err := process.NewProcess(pid)
		if err == nil {
			cmd_line, _ := proc.Cmdline()
			for _, proc_arg := range to_watch {
				want_process := false
				if strings.Contains(cmd_line, proc_arg) {
					want_process = true
					for _, skip_arg := range to_skip {
						if strings.Contains(cmd_line, skip_arg) {
							want_process = false
						}
					}
				}
				if want_process {
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))] = prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Namespace: namespace,
							Name:      "process",
							Help:      "metrics of the process",
						},
						[]string{
							"process",
							"pid",
							"metric",
						},
					)
					// See http://godoc.org/github.com/shirou/gopsutil/process
					// {"rss":2514944,"vms":110858240,"shared":2113536,"text":897024,"lib":0,"data":0,"dirty":36003840}
					// {"cpu":"cpu","user":0.0,"system":0.0,"idle":0.0,"nice":0.0,"iowait":0.0,"irq":0.0,"softirq":0.0,"steal":0.0,"guest":0.0,"guestNice":0.0,"stolen":0.0}
					meminfo, _ := proc.MemoryInfoEx()
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "memory_rss"}).Set(float64(meminfo.RSS))
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "memory_vms"}).Set(float64(meminfo.VMS))
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "memory_shared"}).Set(float64(meminfo.Shared))
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "memory_text"}).Set(float64(meminfo.Text))
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "memory_lib"}).Set(float64(meminfo.Lib))
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "memory_data"}).Set(float64(meminfo.Data))
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "memory_dirty"}).Set(float64(meminfo.Dirty))
					cpuinfo, _ := proc.Times()
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_user"}).Set(cpuinfo.User)
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_system"}).Set(cpuinfo.System)
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_idle"}).Set(cpuinfo.Idle)
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_nice"}).Set(cpuinfo.Nice)
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_iowait"}).Set(cpuinfo.Iowait)
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_irq"}).Set(cpuinfo.Irq)
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_softirq"}).Set(cpuinfo.Softirq)
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_steal"}).Set(cpuinfo.Steal)
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_guest"}).Set(cpuinfo.Guest)
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_guestNice"}).Set(cpuinfo.GuestNice)
					e.processMetrics[proc_arg+strconv.Itoa(int(pid))].With(prometheus.Labels{"process": proc_arg, "pid": strconv.Itoa(int(pid)), "metric": "cpu_stolen"}).Set(cpuinfo.Stolen)
				}
			}

			//if err == nil {
			//	meminfo, err := proc.MemoryInfoEx()
			//	fmt.Println(meminfo)
			//	if err == nil {
			//		fmt.Println("Memory info :")
			//		fmt.Println("RSS: " + strconv.FormatUint(meminfo.RSS, 10))
			//	}
			//	fmt.Println(proc.Times())
			//}
		}
	}

	e.up.Set(1)
}

func (e *Exporter) resetMetrics() {
	for _, m := range e.processMetrics {
		m.Reset()
	}
}

func (e *Exporter) collectMetrics(metrics chan<- prometheus.Metric) {
	for _, m := range e.processMetrics {
		m.Collect(metrics)
	}
}

// Describe describes all the metrics ever exported by the HAProxy exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up.Desc()
	ch <- e.totalScrapes.Desc()
}

// Collect fetches the stats from configured HAProxy location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.resetMetrics()
	e.scrape()

	ch <- e.up
	ch <- e.totalScrapes
	e.collectMetrics(ch)
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
		processExpr     = flag.String("process.watch", "", processHelpText)
		processExprSkip = flag.String("process.nowatch", "", processHelpText)
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("process_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting process_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	log.Infoln("Watch:", *processExpr)
	log.Infoln("Skip:", *processExprSkip)

	to_watch = strings.Split(*processExpr, ",")
	to_skip = strings.Split(*processExprSkip, ",")

	metrics := map[string]*prometheus.GaugeVec{}
	exporter, err := NewExporter(metrics)
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("process_exporter"))

	log.Infoln("Listening on", *listenAddress)
	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>Process Exporter</title></head>
		<body>
		<h1>Process Exporter</h1>
		<p><a href='` + *metricsPath + `'>Metrics</a></p>
		</body>
		</html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
