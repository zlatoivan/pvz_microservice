package metric

import "github.com/prometheus/client_golang/prometheus"

var (
	GivenOutOrdersCounterMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "given_out_orders_counter",
		Help: "Total number of given out orders."},
	)

	ClientGivenOutOrdersCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "client_given_out_orders_counter",
		Help: "Total number of client given out orders."},
		[]string{"client_id"},
	)

	ReturnedOrdersCounterMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "returned_orders_counter",
		Help: "Total number of returned orders."},
	)

	DeletedPVZsCounterMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "deleted_pvzs_counter",
		Help: "Total number of deleted PVZs."},
	)
)
