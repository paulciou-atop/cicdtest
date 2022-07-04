package main

import (
	"exporter/api/v1/scan"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	devices := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "snmpscan",
		Subsystem: "blob_storage",
		Name:      "snmp_devices",
		Help:      "Number of atop devices found in the LAN.",
	})
	prometheus.MustRegister(devices)
	go func() {
		for {
			devinfos, err := scan.Scan("192.168.13.1/24")
			if err != nil {
				log.Info("scan.Scan error err", err)
				continue
			}
			time.Sleep(30 * time.Second)
			devices.Set(float64(len(devinfos)))
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
