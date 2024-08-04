package pool

import "sync"

// There are serveral roles: Job, Worker, WorkerPool

/**************************Job**************************/
// Job
type Job interface {
	RunTask(request interface{})
}

type JobFunc func()

func (j JobFunc) RunTask(_ interface{}) {
	j()
}

// A channel to take all the taskes
type JobChan chan Job

/*************************Worker*************************/
type Worker struct {
	JobQueue JobChan

	//exist sign
	Quit chan bool

	wg *sync.WaitGroup
}

// Create a new worker
func NewWorker(wg *sync.WaitGroup) Worker {
	return Worker{
		JobQueue: make(JobChan),
		Quit:     make(chan bool),
		wg:       wg,
	}
}

// Start a worker to listen to Job
// When a worker finishes a job, worker need to be re-sent to the pool
func (w Worker) Start(workerPool *WorkerPool) {
	w.wg.Add(1)
	// Create a new goroutine to avoid block
	go func() {
		defer w.wg.Done()
		for {
			//Register the worker to pool
			workerPool.workerQueue <- &w
			select {
			case job := <-w.JobQueue:
				job.RunTask(nil)

			case <-w.Quit:
				return
			}
		}
	}()
}

/**************************Pool**************************/
type WorkerPool struct {
	Size        int
	JobQueue    JobChan
	workerQueue chan *Worker
	wg          sync.WaitGroup
}

// Create a new workerPool
func NewWorkerPool(poolSize, jobQueueLen int) *WorkerPool {
	return &WorkerPool{
		Size:        poolSize,
		JobQueue:    make(JobChan, jobQueueLen),
		workerQueue: make(chan *Worker, poolSize),
	}
}

func (wp *WorkerPool) Start() {

	// Start all workers
	for i := 0; i < wp.Size; i++ {
		worker := NewWorker(&wp.wg)
		worker.Start(wp)
	}

	// Listen to JobQueue, if request is received, randomly pick a worker and send the job to the JobQueue of the worker
	// Need a new goroutine to avoid blocking
	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				worker := <-wp.workerQueue
				worker.JobQueue <- job
			}
		}
	}()
}

func (wp *WorkerPool) Submit(job Job) {
	wp.JobQueue <- job
}

func (wp *WorkerPool) Shutdown() {
	for i := 0; i < wp.Size; i++ {
		worker := <-wp.workerQueue // Get waitting workers
		worker.Quit <- true        // Send exist signal
	}
	wp.wg.Wait() // Waitting for all workers done
}
