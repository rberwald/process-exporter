# process-exporter


Software for monitoring of individual processes on Linux with Prometheus.

### Description

This repository is the home of the process-exporter, a prometheus exporter to monitor individual processes on Linux.

The exporter takes two parameters to select and de-select which processes to monitor.

### Arguments

```Usage of /process-exporter:
  -log.format value
    	Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true" (default "logger:stderr")
  -log.level value
    	Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
  -process.nowatch string
    	Processes to (no)watch
	You should provided at least one process to watch.
	The parameter process.watch should be a comma-seperated list of regular expressions of processes to watch
	The parameter process.nowatch is a filter that removes processes from the list provided by process.watch
  -process.watch string
    	Processes to (no)watch
	You should provided at least one process to watch.
	The parameter process.watch should be a comma-seperated list of regular expressions of processes to watch
	The parameter process.nowatch is a filter that removes processes from the list provided by process.watch
  -version
    	Print version information.
  -web.listen-address string
    	Address to listen on for web interface and telemetry. (default ":8980")
  -web.telemetry-path string
    	Path under which to expose metrics. (default "/metrics")
```

### Usage

The selection of the processes to return is rather simple.

* Of every process found, the command line is taken.
* The contents of the command line is matched with the arguments of process.watch. No match, go to next proccess.
* Next, match the command line with the arguments of process.nowatch, Match, go to next process
* If process should be watch, add a record for this process, labeling it with the argument from process.watch and the pid found.

This way of selecting has some side effects:

- the process of the exporter itself, has all arguments of process.watch in its command line
- If you run 'curl http://localhost:8980/metrics | grep \<process\>' where process is in the process.watch list, you will get an extra record in the output for each metric, since 'grep \<process\>' actually matches.

It's up to the user to add enough arguments to process.nowatch to make sure only the required processes are being reporter. It's advised to at least give process-expoter as an argument to process.nowatch.


#### Running process-exporter inside docker

If you want to run the process-exporter inside docker, you need to make sure the process-exporter has access to the process information of the host. This can be done by mounting the /proc filesystem to the docker container, and setting an environment variable to point the process-exporter to the right directory.

Docker example:

```
docker run --rm -d -v /:/rootfs -e HOST_PROC=/rootfs/proc -e HOST_SYS=/rootfs/sys -e HOST_ETC=/rootfs/etc -p 8980:8980 --name=process-exporter rberwald/process-exporter /process-exporter -process.watch process1,process2 -process.nowatch process-exporter
```

#### Usage in Prometheus

sum by (process) (irate(indivdual_process{instance="\<server\>:8980", metric=~"cpu_.*"}[1m]))

indivdual_process{instance="\<server\>:8980", metric=~"memory_.*", process="\<process\>"}


### Environment variables

The process-exporter itself is not using any environment variables.

The library used to get the statistics of the processes (https://github.com/shirou/gopsutil/tree/master/process) is using these environment variables:
- HOST_PROC
- HOST_SYS
- HOST_ETC

Setting these variables is necessary when running the process-exporter inside a docker containers.

### Metrics

Overview of the metrics exporter by this exporter.

indivdual_process{metric="cpu_guest",pid="915",process="\<process\>"}
indivdual_process{metric="cpu_guestNice",pid="915",process="\<process\>"}
indivdual_process{metric="cpu_idle",pid="915",process="\<process\>"}
indivdual_process{metric="cpu_iowait",pid="915",process="\<process\>"}
indivdual_process{metric="cpu_irq",pid="915",process="\<process\>"}
indivdual_process{metric="cpu_nice",pid="915",process="\<process\>"}
indivdual_process{metric="cpu_softirq",pid="915",process="\<process\>"}
indivdual_process{metric="cpu_steal",pid="915",process="\<process\>"}
indivdual_process{metric="cpu_stolen",pid="915",process="\<process\>"}
indivdual_process{metric="cpu_system",pid="915",process="\<process\>"}
indivdual_process{metric="cpu_user",pid="915",process="\<process\>"}
indivdual_process{metric="memory_data",pid="915",process="\<process\>"}
indivdual_process{metric="memory_dirty",pid="915",process="\<process\>"}
indivdual_process{metric="memory_lib",pid="915",process="\<process\>"}
indivdual_process{metric="memory_rss",pid="915",process="\<process\>"}
indivdual_process{metric="memory_shared",pid="915",process="\<process\>"}
indivdual_process{metric="memory_text",pid="915",process="\<process\>"}
indivdual_process{metric="memory_vms",pid="915",process="\<process\>"}

Where process is the name as given in the arguments process.watch.

### Support Scripts
There are some support scripts in the ```scripts/``` sub-directory. The scripts:
- ```build``` : Script to compile the GO code in the ```golang``` container and build a Docker container with the binary.
- ```run``` : Script to quickly run the result of ```build``` on a local docker instance.

### To Do

List of things to do in the future versions:

- Make the selection of processes smarter
- Solve all inline ToDo remarks
- Add automated tests
- Add automated build

