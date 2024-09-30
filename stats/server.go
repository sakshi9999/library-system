package stats

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var bootTime float64

var (
	upTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "uptime",
		Help: "Uptime of service"},
		[]string{"service_name"})
	AddBookApiCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "add_book_api_count",
		Help: "Counter for add book api",
	})
	GetBookApiCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "get_book_api_count",
		Help: "Counter for get book api",
	})
	BorrowBookApiCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "borrow_book_api_count",
		Help: "Counter for borrow book api",
	})
	ReturnBookApiCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "return_book_api_count",
		Help: "Counter for return book api",
	})
	ApiElapsedTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "api_elapsed_time",
		Help: "Time taken to execute an API",
	}, []string{"api_name"})
)

func init() {
	prometheus.MustRegister(upTime)
	prometheus.MustRegister(AddBookApiCounter)
	prometheus.MustRegister(GetBookApiCounter)
	prometheus.MustRegister(BorrowBookApiCounter)
	prometheus.MustRegister(ReturnBookApiCounter)
	prometheus.MustRegister(ApiElapsedTime)

}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func LaunchMetricObserver(wg *sync.WaitGroup, observerAddr string, serverName string) {
	defer wg.Done()
	wg.Add(1)
	go func(serverName string) {
		defer wg.Done()
		for {
			select {
			case <-time.After(10 * time.Second):
				bootTime += 10.0
				upTime.WithLabelValues(serverName).Set(bootTime)
			}

		}
	}(serverName)
	upTime.WithLabelValues(serverName).Set(bootTime)
	r := gin.New()
	r.GET("/metrics", prometheusHandler())
	r.Run(observerAddr)
}
