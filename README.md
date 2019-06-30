# NUT Exporter - WIP

This is an NUT (Network UPS Tools) exporter for Prometheus.

The exporter relies on tools from the https://github.com/networkupstools/nut (uspc).

## Running

A minimal invocation looks like this:

    ./nut_exporter

Supported parameters include:

 - `port`: port to listen on (default: `8100`)
 - `ups`: name of the monitored ups  (default: `none`)
 - `upsc`: path to the upsc executable (default: rely on `$PATH`)
 
 ## Running in docker
    docker run -d -p 8100:8100 -v /bin/upsc:/bin/upsc:ro quay.io/klippo/nut_exporter -ups <upsname>
