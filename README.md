# Sphinx Exporter [![Build Status](https://travis-ci.org/iamseth/sphinx_exporter.svg)](https://travis-ci.org/iamseth/sphinx_exporter.svg)


Export [Sphinx Search](http://sphinxsearch.com/) metrics to Prometheus.

## Exported Metrics

|Metric|Meaning|
|------|------|
|sphinx_agent_connect||
|sphinx_agent_retry||
|sphinx_command_delete|Sum of delete commands since startup.|
|sphinx_command_excerpt|Sum of excerpt commands since startup.|
|sphinx_command_flushattrs|Sum of flushattrs commands since startup.|
|sphinx_command_keywords|Sum of keywords commands since startup.|
|sphinx_command_persist|Sum of persist commands since startup.|
|sphinx_command_search|Sum of search commands since startup.|
|sphinx_command_status|Sum of status commands since startup.|
|sphinx_command_update|Sum of update commands since startup.|
|sphinx_connections|Total number of connections made since startup.|
|sphinx_dist_local||
|sphinx_dist_queries||
|sphinx_dist_wait||
|sphinx_dist_wall|Total time running distributed queries.|
|sphinx_maxed_out|Number of times connections were maxed out.|
|sphinx_queries|Total number of queries run against Sphinx.|
|sphinx_query_cpu||
|sphinx_query_readkb||
|sphinx_query_reads||
|sphinx_query_readtime||
|sphinx_query_wall|Total time running queries.|
|sphinx_up|1 if we're able to scrape metrics, otherwise 0.|
|sphinx_uptime|Time in seconds searchd has been running.|

## Flags

```bash
./sphinx_exporter --help
```

* __`sphinx.host`:__ Hostname or IP address for Sphinx. (default "localhost")
* __`sphinx.port`:__ TCP port for Sphinx.
* __`web.listen-addr`:__ Address to listen on for web interface and telemetry. (default "0.0.0.0:9247")
* __`web.telemetry-path`:__ Path under which to expose metrics. (default "/metrics")

## Useful Queries

TODO

## Using Docker

You can deploy this exporter using the [iamseth/sphinx_exporter](https://registry.hub.docker.com/u/iamseth/sphinx_exporter) Docker image.

For example:

```bash
docker pull iamseth/sphinx_exporter

docker run -d -p 9247:9247 iamseth/sphinx_exporter -sphinx.host 192.168.1.100
```

## Installation

```bash
sudo useradd sphinx_exporter
sudo curl -Ls "https://raw.githubusercontent.com/iamseth/sphinx_exporter/master/examples/systemd/sphinx_exporter.service" -o /etc/systemd/system/sphinx_exporter.service
sudo curl -Ls "https://raw.githubusercontent.com/iamseth/sphinx_exporter/master/examples/systemd/sysconfig.sphinx_exporter" -o /etc/sysconfig/sphinx_exporter
sudo curl -Ls "https://github.com/iamseth/sphinx_exporter/releases/download/0.0.2/sphinx_exporter.linux-amd64" -o /usr/sbin/sphinx_exporter
sudo chmod 755 /usr/sbin/sphinx_exporter
sudo sytemctl start sphinx_exporter
```
