package main

import (
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var onlineUsers = prometheus.NewGauge(prometheus.GaugeOpts{
  Name: "goapp_online_users",
  Help: "Online users",
  ConstLabels: map[string]string{
    "course": "full cycle",
  },
})

var httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
  Name: "goapp_http_requests_total",
  Help: "count of all http requests for goapp",
}, []string{})

var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
  Name: "goapp_http_requests_duration",
  Help: "duration in seconds of all http requests",
}, []string{"handler"})

func produceVariableOnlineUsers() {
  for{
    onlineUsers.Set(float64(rand.Intn(2000)))
  }
}

func main() {
  register := prometheus.NewRegistry()
  register.MustRegister(onlineUsers)
  register.MustRegister(httpRequestsTotal)
  register.MustRegister(httpDuration)

  go produceVariableOnlineUsers()

  http.Handle("/metrics", promhttp.HandlerFor(register ,promhttp.HandlerOpts{}))

  home := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("hello full cycle\n"))
  })

  duration := promhttp.InstrumentHandlerDuration(
    httpDuration.MustCurryWith(prometheus.Labels{"handler":"home"}),
    promhttp.InstrumentHandlerCounter(httpRequestsTotal, home),
  )

  http.Handle("/", promhttp.InstrumentHandlerCounter(httpRequestsTotal, duration))
  http.ListenAndServe(":8181", nil)
}