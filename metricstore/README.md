## Metric data store

This is different then the Data store service which stores general configuration and user data in a transactional way.

This data store is for massive bulk-data storage, metric data, telemetry and events and logs.

It is targeting time-series databases like influxdb, Prometheus, Thanos, druid, cortex and m3db, etc.

Make it not dependent on a specific database. Create API layer and insulation to allow for different backends.

