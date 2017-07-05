# Can't use scratch or Alpine, since net/http is not static :(
FROM prom/busybox:glibc

COPY process-exporter /process-exporter

EXPOSE 8980
