statsd-firehose
---------------

`statsd-firehose` is a simple load-testing tool intended for stressing statsd
and, more importantly, whatever you have behind statsd. In my case, that's
carbon-cache, whisper, and network-attached SSDs.


    Usage of ./statsd-firehose:
      -statsd="statsd://127.0.0.1:8125/firehose.": Statsd URL including a prefix for all metrics
      -packetsize=512: UDP packet size for metrics sent to statsd
      -countcount=50000: Number of individual counters to run
      -countinterval=60: Gauge update interval, in seconds
      -gaugecount=50000: Number of individual gauges to run
      -gaugeinterval=60: Gauge update interval, in seconds

