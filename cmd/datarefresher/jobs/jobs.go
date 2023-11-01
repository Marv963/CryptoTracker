package jobs

import (
	"fmt"
	"time"

	"github.com/Marv963/CryptoTracker/app/pkg/appcontext"
	"github.com/Marv963/CryptoTracker/app/pkg/workerpool"
)

type JobController struct {
	appContext *appcontext.AppContext
	pool       *workerpool.WorkerPool
	wsURL      string
}

func NewJobController(appContext *appcontext.AppContext, wsURL string) *JobController {
	return &JobController{
		appContext: appContext,
		pool:       workerpool.NewWorkerPool(4),
		wsURL:      wsURL,
	}
}

func (s *JobController) FetchDataWithInterval(job func(), interval time.Duration) {
	for {
		start := time.Now() // Time of the start of the data retrieval
		s.appContext.Logger.Println("Fetching data at", time.Now())
		job()
		elapsed := time.Since(start) // Dauer des Datenabrufs
		s.appContext.Logger.Println("Finishing Fetching data at", time.Now(), " took", elapsed)

		// If the duration of the data retrieval is less than the desired interval,
		// wait the rest of the time to complete the interval.
		if elapsed < interval {
			s.appContext.Logger.Println("Waiting for", interval-elapsed)
			time.Sleep(interval - elapsed)
		}
	}
}

func (s *JobController) Shutdown() {
	// Stop the WorkerPool
	s.pool.Stop()
}

func (s *JobController) Start() {
	fmt.Println("JobController started")
	go s.handleWorkerErrors()
}

// New error-handling method/goroutine
func (s *JobController) handleWorkerErrors() {
	for task := range s.pool.Errors {
		fmt.Printf("Handling error for pair: %s err: %v \n", task.Pair, task.Error)
	}
}
