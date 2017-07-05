# rooomy-exporter

Software for blackbox monitoring of RoOomy instances with Prometheus.

### Description

This repository is the home of the rooomy-exporter, a prometheus exporter to monitor RoOomy instances.

The exporter takes an URI as argument to point itself to the RoOomy instances. The URI needs to contain the username and password of the RoOomy instances.

Every time the metrics are gathered, the exporter will contact the RoOomy instances, retrieve the AYT information, and use the AYT information to retrieve the Gallery, Real Estate Property, Model Product and User Profile feed. It will report the number of milliseconds of retrieval as metric.

### Arguments

```Usage of /rooomy-exporter:
  -log.format value
    	Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true" (default "logger:stderr")
  -log.level value
    	Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
  -rooomy.scrape-uri string
    	URI on which to scrape rooomy. (default "http://localhost/")
  -version
    	Print version information.
  -web.listen-address string
    	Address to listen on for web interface and telemetry. (default ":8989")
  -web.telemetry-path string
    	Path under which to expose metrics. (default "/metrics")
```

### Environment variables

- ROOOMY_USER : the user to use to access the rooomy instance
- ROOOMY_PASSWORD : the password to use to access the rooomy instance
- ROOOMY_ENVIRONMENT : the environment of the rooomy instance (eg test)
- ROOOMY_ROLE : the role of the rooomy instance (slave, master)

The variables ROOOMY_USER and ROOOMY_PASSWORD, if set to non-empty, overwrite credentials of the URI.

### Metrics

Overview of the metrics exporter by this exporter.

- rooomy_exporter_build_info : version, build and source code information
- rooomy_lengthAYT : length of AYT response
- rooomy_timeAYT : Time it took to get AYT information, in milliseconds
- rooomy_timeGallery : Time it took to get Gallery Feed, in milliseconds
- rooomy_timeModelProduct : Time it took to get Model Product Feed, in milliseconds
- rooomy_timeRealEstate : Time it took to get Real Estate Property Feed, in milliseconds
- rooomy_timeUserProfile : Time it took to get User Profile, in milliseconds
- rooomy_countApiRequest : Number of API requests
- rooomy_countHttp200 : Number of requests leading to HTTP 2xx response
- rooomy_countHttp400 : Number of requests leading to HTTP 4xx response
- rooomy_countHttp500 : Number of requests leading to HTTP 5xx response
- rooomy_countSqlConnections : Number of available SQL connections
- rooomy_countSqlConnectionsBusy : Number of busy SQL connections
- rooomy_countTcpConnectionsAccepted : Number of accepted TCP connections
- rooomy_countTcpConnectionsActive : Number of active TCP connections
- rooomy_countTcpServerThreads : Number of TCP server threads
- rooomy_countWorkerThreads : Number of available Worker threads
- rooomy_countWorkerThreadsBusy : Number of busy Worker threads
- rooomy_ApplicationInfo : Labelling info (application_name, application_version, application_environment, application_role)

### Redirection

The redirection is useful if you are running in a container and want to connect to the rooomy instance running on the host.
Normally you would need to connect to the public IP of the rooomy instance, sending you through the security groups. This means that the public IP of the host needs to be in the security group of of the host. In Terraform, this leads to a circular definition (sg -> inst -> sg). Plus it adds undesirable delay of going through several firewalls.

By defining and extra host to the container (--add-host="rooomyhost.loftweb.nl:<private ip>"), and using the magical host (rooomyhost.loftweb.nl), you will always stay on the host itself.
You will need determine the private IP of the host to inject in the start command of the container.

### Support Scripts
There are some support scripts in the ```scripts/``` sub-directory. The scripts:
- ```build``` : Script to compile the GO code in the ```golang``` container and build a Docker container with the binary.
- ```run``` : Script to quickly run the result of ```build``` on a local docker instance.

### To Do

List of things to do in the future versions:

- Solve all inline ToDo remarks
- Add automated tests
- Add automated build

