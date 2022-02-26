[![Release Go project](https://github.com/theztd/sentry-exporter/actions/workflows/release.yml/badge.svg)](https://github.com/theztd/sentry-exporter/actions/workflows/release.yml)

# Sentry Exporter for prometheus

Parse sentry stats API and provide metrics to prometheus.

There are two available modes:

**DAEMON mode**

Start webserver providing HTTP access to gathered metrics. Metrics are stored in file on disk and is regenerated every 60s (by default).

This mode will be enabled after **-daemon-mode** parameter is set. 
This mode uses parameters:
 * **-i \<int\>**   ...interval in seconds between reading stats from sparkpost (Example 120)
 * **-p \<string\>** ...http endpoint listening port (Example :9091)
 * **-read-token=\<string\>** ...token that have to be sent by prometheus to read metrics (Example -read-token=MySECRET123 mean that request should looks like /_metrics?token=MySECRET123)


**SIMPLE mode**

Get metrics and generate file in promfile format to defined path and exits. It works well with node_exporter and defined metrics dir.

This mode ignores parameters:
 * -i
 * -p
 * -read-token 


## Metrics example

```
# HELP sentry_received_events_count Hourly received amount of events grouped by project-slug.
sentry_received_events_count{project="fe-landing"} 6
sentry_received_events_count{project="fe-user_area"} 1
```


## Program help
```bash
  -daemon-mode
        Run in daemon mode and server results via HTTP protocol
  -i int
        Interval in seconds between gathering data from sentry (default 60)
  -o string
        Path to file where metrics should be saved (default "tmp_out.prom")
  -p string
        Listening port (default ":9091")
  -read-token string
        Read token for HTTP endpoint, if empty, no autentication is required
```

## Program accept ENV vars

* **SENTRY_TOKEN**   - required parameter
* **INTERVAL** - same to -i parameter
* **PORT** - same to -p parameter
* **METRICS_FILE** - same to -o parameter
* **READ_TOKEN** - same to -read-token parameter