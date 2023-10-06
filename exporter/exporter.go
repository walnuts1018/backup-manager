package exporter

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/walnuts1018/backup-manager/config"
)

var (
	runningTasks = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "running_backup_job_count",
		},
		[]string{},
	)
	jobs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "last_backup_job_success",
		},
		[]string{"name"},
	)
)

func Export() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+config.Config.ServerPort, nil)
}

func SetRunningTasks(tasks int) {
	runningTasks.WithLabelValues().Set(float64(tasks))
}

func SetJobs(name string, lastSuccess bool) {
	if lastSuccess {
		jobs.WithLabelValues(name).Set(1)
	} else {
		jobs.WithLabelValues(name).Set(0)
	}
}
