package queue

// Dispatcher is responsible for managing the workers and the queue
type Dispatcher struct {
	workerCount int
	queue       chan IJob
	workers     []*worker
	managing    chan IJob
}

// GetDispatcher creates and returns a new Dispatcher
func GetDispatcher() *Dispatcher {
	queue := make(chan IJob)
	managing := make(chan IJob)

	managerObj := manager{
		queue:    queue,
		managing: managing}

	go managerObj.start()

	return &Dispatcher{
		queue:    queue,
		managing: managing}
}

// StopAllWorkers actually does what it says it does.
// Note: The workers will not be killed instantly, they will
// finish executing the current job before stopping
func (d *Dispatcher) StopAllWorkers() {
	for _ = range d.workers {
		d.queue <- nil
	}
}

// SetWorkerCount sets a new worker count, and starts all workers,
// if the number of workers is reduced, the workers will finish
// the current job before being stopped
func (d *Dispatcher) SetWorkerCount(n int) {
	if n > 0 {
		d.workerCount = n

		for d.workerCount != len(d.workers) {
			if d.workerCount > len(d.workers) {
				d.addWorker()
			} else {
				d.removeWorker()
			}
		}
	}
}

func (d *Dispatcher) addWorker() {
	worker := newWorker(d.queue)
	go worker.start()
	d.workers = append(d.workers, worker)
}

func (d *Dispatcher) removeWorker() {
	if len(d.workers) > 1 {
		d.queue <- nil
	}
}

// AddJob adds a job to the queue.
// If the dispatcher has no workers, AddJob will create one.
func (d *Dispatcher) AddJob(job IJob) {
	if len(d.workers) == 0 {
		d.addWorker()
	}

	d.managing <- job
}

type manager struct {
	queue    chan IJob
	managing chan IJob
}

func (m *manager) start() {
	for {
		job := <-m.managing
		m.queue <- job

	}
}
