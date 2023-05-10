# weatherflow_exporter

weatherflow_exporter is a Prometheus exporter for
[WeatherFlow Tempest](https://weatherflow.com/tempest-home-weather-system/)
weather stations.

Long-lived WebSocket connections to the
[WeatherFlow Tempest API](https://weatherflow.github.io/Tempest/) are used to
expose metrics from both `obs_st` messages (every 60 sec) and `rapid_wind`
messages (every 3 sec).

Many stations may be monitored simultaneously.  Upon receiving a HTTP request
to the `/scrape` endpoint, a WebSocket connection will be created dynamically
to WeatherFlow, using the `token` and `device_id` passed in as query parameters.
Multiple devices sharing the same token (i.e. same owner) will be multiplexed
over the same connection.  The connection will automatically timeout after 30
minutes of inactivity (i.e. no requests to `/scrape` for its metrics).

## Install

Download from [releases](https://github.com/tris/weatherflow_exporter/releases)
or run from Docker:

```
docker run -d -p 6969:6969 ghcr.io/tris/weatherflow_exporter
```

An alternate port may be defined using the `PORT` environment variable.  There
are no other configuration options.

## Grafana dashboard

TBD.  You can get some inspiration from the
[Mussel Rock Weather](https://mr.ethereal.net) dashboard, which this powers.

[Dave Schmid](https://github.com/lux4rd0) has also built some excellent
dashboards as part of his own
[WeatherFlow Dashboards AIO](https://github.com/lux4rd0/weatherflow-dashboards-aio)
project.  Note that some metrics (e.g. RSSI) are only available via local UDP
broadcasts -- if you need those, consider using his
[weatherflow-collector](https://github.com/lux4rd0/weatherflow-collector) (for
InfluxDB) instead of this one.

## Test

```
curl -s http://localhost:6969/scrape?token=...&device_id=...
```

The first request will return an empty response, as it simply establishes the
WebSocket connection.  Wait 3 seconds before trying again, and you should
receive the full set of metrics.

## Example Prometheus config

```yaml
scrape_configs:

- job_name: 'weatherflow'
  scrape_interval: 3s
  static_configs:
  - targets:
    - 123456
    labels:
      station_name: 'My Station'
      __param_token: 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'
  - targets: 
    - 234567
    labels:
      station_name: 'Other Station'
      __param_token: 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'
  metrics_path: /scrape
  relabel_configs:
  - source_labels: [__address__]
    target_label: __param_device_id
  - target_label: __address__
    # IP of the exporter
    replacement: localhost:6969

- job_name: 'weatherflow_exporter'
  static_configs:
  - targets: ['localhost:6969']
```

## See also

- [weatherflow](https://github.com/tris/weatherflow) Go module
- [WeatherFlow Tempest API](https://weatherflow.github.io/Tempest/)
- [WeatherFlow Tempest API WebSocket Reference](https://weatherflow.github.io/Tempest/api/ws.html)
- [Tempest API Remote Data Access Policy](https://weatherflow.github.io/Tempest/api/remote-developer-policy.html)
- [tempest-exporter](https://github.com/nalbury/tempest-exporter)
- [weatherflow-collector](https://github.com/lux4rd0/weatherflow-collector)
