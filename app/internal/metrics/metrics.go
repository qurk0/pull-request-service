package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	PRCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pr_created_total",
		Help: "Сколько пулл реквестов создано",
	})

	PRMerged = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pr_merged_total",
		Help: "Сколько пулл реквестов замёржено",
	})

	PRReassigned = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pr_reassigned_total",
		Help: "Число успешных переназначений ревьюверов",
	})

	PRReassignFailed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "pr_reassign_failed",
		Help: "Число провальных попыток переназначений ревьюверов",
	}, []string{"reason"})

	TeamsCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "team_created_total",
		Help: "Сколько команд зарегистрировано в сервисе",
	})
)

func Init() {
	prometheus.MustRegister(PRCreated, PRMerged, PRReassigned, PRReassignFailed)
}
