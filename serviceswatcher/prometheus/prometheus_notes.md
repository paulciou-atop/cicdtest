# Alarm service

## Prometheus

### Install
```shell
$ wget https://github.com/prometheus/prometheus/releases/download/v2.35.0/prometheus-2.35.0.linux-amd64.tar.gz
$ tar xvfz prometheus-2.35.0.linux-amd64.tar.gz
```

#### Run
**Run native**
```shell
./prometheus --config.file=prometheus.yml
```
**Node Exporter**
Install  
```shell
https://github.com/prometheus/node_exporter/releases/download/v1.3.1/node_exporter-1.3.1.linux-amd64.tar.gz
```

### Query
https://prometheus.io/docs/prometheus/latest/querying/basics/

### Concept
- Pushgateway: metrics cache, some jobs may not exist long enough to be scraped, they can intead push their metrics to a Pushgateway.
- Alertmanager: handle alerts.


#### Metrics
- Metrics consist of metric name + optional labels (key-value pairs)
- format: `/[a-zA-Z_:][a-zA-Z0-9+:]*/`
- where Prometheus get metrics: `/metrics` http endpoint
  
**Naming convention**
- Should have single word application prefix, this prefix is referred to as a `namespace` by client library
- Shoudl have a suffix describing the unit, in plural form. ex:
  - http_request_duration_**seconds**
  - node_memory_usage_**bytes**

**Format**
`<metric name>{<label name>=<label value>, ...}`
example:  
```yaml
api_http_requests_total{method="POST", handler="/messages"}
```


Real sample of metrics
- HELP: Description of what the metrics is
- TYPE: 
  - Counter: ...how many times x happened
  - Gauge: ...what is the current value of x now
  - Histogram: ...how long or how big 
  - Summary: 
    - φ-quantiles: percentile for example quantile=0.25 means scraped value < 25%. Refer a sample following, `go_gc_duration_seconds{quantile="0.25"} 5.3437e-05` meansgolang gc duration = 5.3437e-05 seconds < 25%
```
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 2.7794e-05
go_gc_duration_seconds{quantile="0.25"} 5.3437e-05
go_gc_duration_seconds{quantile="0.5"} 8.5913e-05
go_gc_duration_seconds{quantile="0.75"} 0.000120427
```

#### Job and instances
- An endpoint can scrape is called an *instance*.
- A collection of *instances* whith same purpose, a process replicated for scalability or reliability is called a *job*.
```yaml
job: api-server
  - instance 1: 1.2.3.4:5670
  - instance 2: 1.2.3.4:5671
```
Every instance scarpe, Prometheus stores a sample in the folliwing time sereis:
- `up{job="<job-name>", instance="<instance-id>"}`: `1` if the instance is healthy, `0` if the scrape failed.
- `scrape_duration_seconds{job="<job-name>", instance="<instance-id>"}` 
- `scrape_samples_post_metric_relabeling{job="<job-name>", instance="<instance-id>"}`
- `scrape_samples_scraped{job="<job-name>", instance="<instance-id>"}`: number of samples the target exposed.
- `scrape_series_added{job="<job-name>", instance="<instance-id>"}`: The approximate number of new series in this scrape.


Prometheus Server
- Storage: stores metrics, Time series database
- Retrieval: A worker pulls metrics data
- Http server: accepts queryies (PromQL)


Alert Manager
- rules file


